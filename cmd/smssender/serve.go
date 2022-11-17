package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kavenegar/kavenegar-go"
	"github.com/spf13/cobra"
)

var (
	smsSender *SmsSender
	serveCmd  = &cobra.Command{
		Use:   "serve",
		Short: "serve smssender",
		Long:  "Serves SMS Sender web app",
		Run:   serve,
	}
)

type SendRequest struct {
	To      string
	From    string
	Message string
	SentAt  time.Time
}

type SendResponse struct {
	status int
}
type Sender interface {
	Send(SendRequest) SendResponse
}

type SmsSender struct {
	config *Config
	api    kavenegar.Kavenegar
}

func serve(cmd *cobra.Command, args []string) {

	config := loadConfigOrPanic(cmd, args)
	smsSender = NewSmsSender(config)

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/send", sendHandler)
	router.HandleFunc("/records", recordsHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.SmsSender.Host, config.SmsSender.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block until channel received os.Interrupt
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), config.SmsSender.GracefulTimeout)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Gracefully shutting down...")
	os.Exit(0)
}

func loadConfigOrPanic(cmd *cobra.Command, args []string) *Config {
	config, err := loadConfig(cmd, args)
	if err != nil {
		panic(err)
	}

	return config
}

func NewSmsSender(config *Config) *SmsSender {
	return &SmsSender{config: config}
}

func (s *SmsSender) Send(req *SendRequest) error {
	resp, err := s.api.Message.Send(
		req.From,
		[]string{req.To},
		req.Message,
		nil,
	)
	if err != nil {
		return err
	}

	// save the record to DB
	go saveToDB(req)

	for _, r := range resp {
		log.Printf("successfully sent message from %s to %s", r.Sender, r.Receptor)
	}
	return nil
}

func saveToDB(req *SendRequest) {
	// TODO: Implement
}
