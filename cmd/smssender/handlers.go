package main

import (
	"io"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the SMS Sender Web Application!")
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

func recordsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}
