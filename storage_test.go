package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountByID(t *testing.T) {
	s, err := NewPostgresStore()
	assert.Nil(t, err)
	res, err := s.GetAccountByID(2)
	assert.Nil(t, err)
	fmt.Println(res)
}
