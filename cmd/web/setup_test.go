package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Exit after m is runned
	os.Exit(m.Run())
}

// myHandler satisfies http handler interface
type myHandler struct{}

// bind methods ServeHTTP to myHandler -> satisfy the handler interface
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Empty method
}
