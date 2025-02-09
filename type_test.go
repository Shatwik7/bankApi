package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	_, err := NewAccount("a", "b", "hunter")
	assert.Nil(t, err)
}
