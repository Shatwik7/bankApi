package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func StoreInContext(r *http.Request, key string, val any) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}
func getIDFromToken(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("INVALID TOKEN")
	}
	claimsID, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("INVALID TOKEN")
	}
	claimsIDInt64 := int64(claimsID)
	return claimsIDInt64, nil
}
func withJWTAuthForTransfer(handerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ApiError{Error: "INVALID TOKEN"})
			return
		}
		userid, err := getIDFromToken(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, err)
			return
		}
		fmt.Println("userId : ", userid)
		r = StoreInContext(r, "userId", userid)
		handerFunc(w, r)
	}
}
func withJWTAuth(handerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ApiError{Error: "INVALID TOKEN"})
			return
		}
		userid, err := getIDFromToken(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, err)
			return
		}
		idstr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: "param should be a integer"})
			return
		}
		fmt.Println("userId : ", id)
		r = StoreInContext(r, "userId", id)
		if userid != int64(id) {
			writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Unuthorized"})
			return
		}
		handerFunc(w, r)
	}
}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "password"
	}
	claims := &jwt.MapClaims{
		"expirestAt":    15000,
		"accountNumber": account.Number,
		"id":            account.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "password"
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error validating")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

type ApiError struct {
	Error string
}
type APIServer struct {
	listenAddr string
	store      Storage
}
type ApiFunc func(http.ResponseWriter, *http.Request) error

func makeHttpHandelFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}

	}
}

func NewAPIserver(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandelFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandelFunc(s.handleAccountById))).Methods(http.MethodGet, http.MethodDelete)
	router.HandleFunc("/transfer", withJWTAuthForTransfer(makeHttpHandelFunc(s.handleTransfer))).Methods(http.MethodPost)
	router.HandleFunc("/login", makeHttpHandelFunc(s.handleLogin)).Methods(http.MethodPost)
	log.Println("BANK API RUNNING ON PORT :", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccounts(w, r)
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}

}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccountById(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccountById(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusAccepted, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	account, err := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)
	if err != nil {
		return writeJSON(w, http.StatusNotAcceptable, ApiError{Error: "Not a Valid Password"})
	}
	log.Println("account:", account)
	if _, err := s.store.CreateAccount(account); err != nil {
		return err
	}
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}
	log.Println("account:", account)
	log.Println("token :", tokenString)
	defer r.Body.Close()
	return writeJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return nil
	}
	res, err := s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusAccepted, res)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginUserRequest := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: "Bad Request"})
	}
	account, err := s.store.GetAccountByID(loginUserRequest.ID)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, err)
	}
	ok := bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(loginUserRequest.Password))
	if ok != nil {
		return writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Wrong"})
	}
	token, err := createJWT(account)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: "Server Error"})
	}
	defer r.Body.Close()
	return writeJSON(w, http.StatusAccepted, &LoginResponse{
		ID:        account.ID,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Token:     token,
		Balance:   account.Balance,
		Number:    account.Number,
	})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TrasnferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	userID, ok := r.Context().Value("userId").(int64)
	if !ok {
		return writeJSON(w, http.StatusBadRequest, ApiError{Error: "failure"})
	}
	s.store.TransferAmount(int(userID), transferReq.ToAccount, transferReq.Amount)
	defer r.Body.Close()

	return writeJSON(w, http.StatusAccepted, transferReq)
}

func getID(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, fmt.Errorf("INVALID ID")
	}
	return id, nil
}
