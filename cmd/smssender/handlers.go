package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang/gddo/httputil/header"
	core "github.com/homayoonalimohammadi/go-sms-sender/smssender/internal/app"
	"github.com/kavenegar/kavenegar-go"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the SMS Sender Web Application!")
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	// check for correct header
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "invalid Content-Type: must be application/json"
			http.Error(w, msg, http.StatusBadRequest)
		}
	}
	// set max bytes to prevent too large request body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)

	var req core.SendRequest
	err := decoder.Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.SentAt = time.Now()
	err = smsSender.Send(&req)
	if err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case *kavenegar.HTTPError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "unknown error from api", http.StatusInternalServerError)
		}
		return
	}

	// save record to postgres
	go func() {
		err = smsSender.postgresProvider.Save(&req)
		if err != nil {
			log.Println("error saving record to postgres:", err)
		}
	}()

	io.WriteString(w, "successfully sent all messages")
}

func RecordsHandler(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")

	sendRequests, err := smsSender.postgresProvider.Get(sender)
	if err != nil {
		log.Println("error retrieving sender records from postgres:", err)
		io.WriteString(w, "Something went wrong, try again shortly...")
		return
	}

	for _, record := range sendRequests {
		io.WriteString(w, fmt.Sprintf("%+v", *record))
	}
}
