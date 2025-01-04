package main

import (
	"log"
	"net/http"

	"github.com/rupert648/todo/adapters/todohttp"
	"github.com/rupert648/todo/adapters/todohttp/views"
	"github.com/rupert648/todo/domain/shopping"
)

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

	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
