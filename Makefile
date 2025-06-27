APP_NAME="bumflix"

build:
	@go build -o bin/$(APP_NAME) src/api/main.go

run: run-web && run-api

run-api:
	@air --build.cmd "go build -o bin/api src/api/main.go" --build.bin "./bin/api"

run-web:
	@cd src/web && pnpm dev

clean:
	@rm -rf bin
	@rm -rf /tmp/bumflix

test:
	@go test -v ./...

start-services:
	@docker compose up -d

# generate password using openssl
# Change password of MINIO_PASSWORD & POSTGRES_PASSWORD in .env file
# docker compose up -d

# compose-up:
# 	@MINIO_PASSWORD=$$(openssl rand -base64 16 | tr -dc 'A-Za-z0-9' | head -c 16) && \
# 		sed -i "s/MINIO_PASSWORD=.*/MINIO_PASSWORD=$$MINIO_PASSWORD/g" .env && \
# 		echo "Generated password for MINIO is $$MINIO_PASSWORD"
# 	@POSTGRES_PASSWORD=$$(openssl rand -base64 16 | tr -dc 'A-Za-z0-9' | head -c 16) && \
# 		sed -i "s/POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$$POSTGRES_PASSWORD/g" .env && \
# 		echo "Generated password for POSTGRES is $$POSTGRES_PASSWORD"
# 	@docker compose up -d
#
# compose-stop:
