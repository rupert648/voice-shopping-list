.PHONY: serve-web

serve-web:
	cd app && go run cmd/server/main.go
