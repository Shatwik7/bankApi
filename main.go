package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("store ", store)
	fmt.Println("HELLO")
	server := NewAPIserver(":3000", store)
	server.Run()
}
