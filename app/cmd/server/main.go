package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rupert648/todo/adapters/todohttp"
	"github.com/rupert648/todo/adapters/todohttp/views"
	"github.com/rupert648/todo/domain/shopping"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf(
			"Started %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}

const addr = ":8000"

func main() {
	list, err := shopping.NewList("shopping.db")
	if err != nil {
		log.Fatal(err)
	}

	templates, err := views.NewTemplates()

	if err != nil {
		log.Fatal(err)
	}

	handler, err := todohttp.NewShoppingHandler(list, views.NewTodoView(templates), views.NewIndexView(templates))

	if err != nil {
		log.Fatal(err)
	}

	loggedHandler := LoggingMiddleware(handler)

	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, loggedHandler); err != nil {
		log.Fatal(err)
	}
}
