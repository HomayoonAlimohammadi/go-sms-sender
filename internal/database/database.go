package database

import (
	"database/sql"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	core "github.com/homayoonalimohammadi/go-sms-sender/smssender/internal/app"
	_ "github.com/lib/pq"
)

type Database interface {
	Save(*core.SendRequest) error
	Get(string) ([]*core.SendRequest, error)
}

type PostgresConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	DbName         string
	SslMode        string
	MigrationsPath string
}

type PostgresProvider struct {
	Config *PostgresConfig
	DB     *sql.DB
}

func (c *PostgresConfig) GetUrl() *url.URL {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Path:   c.Host + ":" + c.Port + "/" + c.DbName,
	}

	query := pgURL.Query()
	query.Add("sslmode", c.SslMode)
	pgURL.RawQuery = query.Encode()
	return pgURL
}

func NewPostgresProvider(config *PostgresConfig) (*PostgresProvider, error) {
	pgURL := config.GetUrl()
	db, err := sql.Open("postgres", pgURL.String())
	if err != nil {
		return nil, err
	}
	return &PostgresProvider{
		Config: config,
		DB:     db,
	}, nil
}

func (p *PostgresProvider) Close() {
	p.DB.Close()
	log.Println("Closed connection to the Postgres Database")
}

func (p *PostgresProvider) Save(req *core.SendRequest) error {
	statement := "INSERT INTO smssender(sender, receptor, message, sent_at) VALUES ($1, $2, $3, $4)"
	_, err := p.DB.Exec(statement, req.Sender, req.Receptor, req.Message, req.SentAt)
	return err
}

func (p *PostgresProvider) Get(sender string) ([]*core.SendRequest, error) {
	query := "SELECT sender, receptor, message, sent_at FROM smssender;"
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sendRequests []*core.SendRequest
	for rows.Next() {
		sendReq := &core.SendRequest{}
		err = rows.Scan(&sendReq.Sender, &sendReq.Receptor, &sendReq.Message, &sendReq.SentAt)
		if err != nil {
			return nil, err
		}
		sendRequests = append(sendRequests, sendReq)
	}
	return sendRequests, nil
}

func (p *PostgresProvider) Migrate() error {
	driver, err := postgres.WithInstance(p.DB, &postgres.Config{
		MultiStatementEnabled: true,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		p.Config.MigrationsPath,
		p.Config.DbName,
		driver,
	)
	if err != nil {
		return err
	}

	var migrationError error
	for migrationError == nil {
		migrationError = m.Steps(1)
	}

	log.Println("Applied migrations successfully.")
	return nil
}
