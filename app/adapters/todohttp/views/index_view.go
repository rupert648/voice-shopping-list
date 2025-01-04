package views

import (
	"html/template"
	"net/http"

	"github.com/rupert648/todo/domain/shopping"
)

type IndexView struct {
	templ *template.Template
}

func NewIndexView(templ *template.Template) *IndexView {
	return &IndexView{templ: templ}
}

func (t *IndexView) Index(w http.ResponseWriter, shoppingItems []shopping.ShoppingItem) {
	var viewModel any = shoppingItems
	if err := t.templ.ExecuteTemplate(w, "index", viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
