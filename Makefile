.PHONY: all
all: server cli

bin:
	mkdir bin

.PHONY: server
server: bin
	go build -o bin/ ./cmd/raccoonstash-server

.PHONY: cli
cli: bin
	go build -o bin/ ./cmd/raccoonstash-cli

.PHONY: sql
sql:
	sqlc generate