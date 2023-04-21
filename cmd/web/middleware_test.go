package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	// create a variable that satisfies http handler's interface
	var myH myHandler

	h := NoSurf(&myH) // pass myH as a pointer
	// test if returned h is a http handler
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	// create a variable that satisfies http handler's interface
	var myH myHandler

	h := SessionLoad(&myH) // pass myH as a pointer
	// test if returned h is a http handler
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
