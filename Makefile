.PHONY: help deps test smssender clean

COVERAGE_PROFILE := c.out 
BINARY_NAME := smssender 


help:
	@echo "This is SMS Sender Web app. Run 'make smssender' and then execute the binary file created to start the app."

deps: 
	go mod download 

test: deps
	go test ./... -coverprofile ${COVERAGE_PROFILE}

build: deps
	go build -o ${BINARY_NAME} ./cmd/smssender

serve: build
	./${BINARY_NAME} serve

clean:
	rm -rf ${BINARY_NAME}
