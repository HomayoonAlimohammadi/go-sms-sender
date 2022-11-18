.PHONY: help deps test smssender clean

COVERAGE_PROFILE := c.out 
BINARY_NAME := smssender 
MIGRATIONS_PATH := internal/database/migrations
POSTGRES_URI := "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}"

help:
	@echo "This is SMS Sender Web app. Run 'make smssender' and then execute the binary file created to start the app."
	@echo Postgres URI: ${POSTGRES_URI}

deps: 
	go mod download 

test: deps
	go test ./... -coverprofile ${COVERAGE_PROFILE}

build: deps clean
	go build -o ${BINARY_NAME} ./cmd/smssender

serve: build
	./${BINARY_NAME} serve

migrate_up:
	migrate -path ${MIGRATIONS_PATH} -database ${POSTGRES_URI} -verbose up

migrate_down:
	migrate -path ${MIGARTIONS_PATH} -database ${POSTGRES_URI} -verbose down

clean:
	rm -rf ${BINARY_NAME}
