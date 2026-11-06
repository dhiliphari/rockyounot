APP_NAME := rockyounot
APP_PATH := ./cmd
SEED_FILE := ./data/rockyou.txt

.PHONY: run build certs up down data

up:
	@echo "Starting Redis..."
	docker compose up -d

down:
	@echo "Stopping Redis..."
	docker compose down

data: $(SEED_FILE)

# Order-only prerequisite: unzip only when the .txt is missing, not when the zip is newer
$(SEED_FILE): | $(SEED_FILE).zip
	@echo "Unzipping $(SEED_FILE).zip..."
	unzip -n $(SEED_FILE).zip -d ./data

run: data
	@echo "Running $(APP_NAME)..."
	go run $(APP_PATH)

build: data
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