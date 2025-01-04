package views

import (
	"html/template"

	"github.com/rupert648/todo/domain/shopping"
)

func NewTodoView(templ *template.Template) *ModelView[shopping.ShoppingItem] {
	return NewModelView[shopping.ShoppingItem](templ, "todo")
}
