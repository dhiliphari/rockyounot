APP_NAME := rockyounot
APP_PATH := ./cmd

.PHONY: run build certs

run:
	@echo "Running $(APP_NAME)..."
	go run $(APP_PATH)

build:
	@echo "Building $(APP_NAME)..."
	mkdir -p ./bin/data
	cp -r ./tls ./bin/
	cp ./data/*.txt ./bin/data/
	go build -o ./bin/$(APP_NAME) $(APP_PATH)

certs:
	@echo "Generating self-signed certificate..."
	mkdir -p tls
	cd tls
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits 2048 --host localhost