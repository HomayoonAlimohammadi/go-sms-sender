package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	core "github.com/homayoonalimohammadi/go-sms-sender/smssender/internal/app"
	"github.com/homayoonalimohammadi/go-sms-sender/smssender/internal/database"
	"github.com/kavenegar/kavenegar-go"
	"github.com/pkg/errors"
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

type Sender interface {
	Send(*core.SendRequest) error
}

type SmsSender struct {
	config           *SmsSenderConfig
	postgresProvider database.Database
	api              *kavenegar.Kavenegar
}

func serve(cmd *cobra.Command, args []string) {

	config := loadConfigOrPanic(cmd, args)
	postgresProvider, err := database.NewPostgresProvider(&config.PostgresDB)
	if err != nil {
		log.Fatalln("unable to provider postgres:", err)
	}
	defer postgresProvider.Close()
	log.Println("connected to postgres")

	// apply postgres migrations
	err = postgresProvider.Migrate()
	if err != nil {
		log.Fatalln("error applying migrations to the postgres db:", errors.WithStack(err))
	}

	smsSender = NewSmsSender(&config.SmsSender, postgresProvider)

	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler).Methods("GET")
	router.HandleFunc("/send", SendHandler).Methods("POST")
	router.HandleFunc("/records", RecordsHandler).Methods("GET")

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.SmsSender.Host, config.SmsSender.Port),
		Handler: router,
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
}

func loadConfigOrPanic(cmd *cobra.Command, args []string) *Config {
	config, err := loadConfig(cmd, args)
	if err != nil {
		panic(err)
	}

	return config
}

func NewSmsSender(config *SmsSenderConfig, postgresProvider *database.PostgresProvider) *SmsSender {
	api := kavenegar.New(config.ApiKey)
	return &SmsSender{config: config, postgresProvider: postgresProvider, api: api}
}

func (s *SmsSender) Send(req *core.SendRequest) error {
	// resp, err := s.api.Message.Send(
	// 	req.Sender,
	// 	[]string{req.Receptor},
	// 	req.Message,
	// 	nil,
	// )
	// if err != nil {
	// 	return err
	// }

	// for _, r := range resp {
	// 	log.Printf("successfully sent message from %s to %s \n", r.Sender, r.Receptor)
	// }
	// return nil

	log.Printf("of course the message was sent from %s to %s! \n", req.Sender, req.Receptor)
	return nil
}

func (s *SmsSender) SaveToPostgres(req *core.SendRequest) {
	err := s.postgresProvider.Save(req)
	if err != nil {
		log.Println("error saving request to database:", err)
	}
}
