.PHONY: build serve-web serve-web-dev

build:
	cd app && go build -o ../bin/server cmd/server/main.go

serve-web: build
	./bin/server

serve-web-dev:
	cd app && go run cmd/server/main.go
