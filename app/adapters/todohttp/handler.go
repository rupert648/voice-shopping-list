package todohttp

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rupert648/todo/adapters/todohttp/views"
	"github.com/rupert648/todo/domain/shopping"
)

var (
	//go:embed static
	static embed.FS
)

type ShoppingHandler struct {
	http.Handler

	list      *shopping.List
	todoView  *views.ModelView[shopping.ShoppingItem]
	indexView *views.IndexView
}

func NewShoppingHandler(service *shopping.List, todoView *views.ModelView[shopping.ShoppingItem], indexView *views.IndexView) (*ShoppingHandler, error) {
	router := mux.NewRouter()
	handler := &ShoppingHandler{
		Handler:   router,
		list:      service,
		todoView:  todoView,
		indexView: indexView,
	}

	staticHandler, err := newStaticHandler()
	if err != nil {
		return nil, fmt.Errorf("problem making static resources handler: %w", err)
	}

	router.HandleFunc("/", handler.index).Methods(http.MethodGet)

	router.HandleFunc("/shopping-item", handler.add).Methods(http.MethodPost)
	router.HandleFunc("/shopping-item", handler.search).Methods(http.MethodGet)
	router.HandleFunc("/shopping-item/sort", handler.reOrder).Methods(http.MethodPost)
	router.HandleFunc("/shopping-item/{ID}/toggle", handler.toggle).Methods(http.MethodPost)
	router.HandleFunc("/shopping-item/{ID}", handler.delete).Methods(http.MethodDelete)
	router.HandleFunc("/shopping-item/{ID}", handler.rename).Methods(http.MethodPatch)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	return handler, nil
}

func (t *ShoppingHandler) index(w http.ResponseWriter, _ *http.Request) {
	items, err := t.list.ShoppingItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.indexView.Index(w, items)
}

func (t *ShoppingHandler) add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := t.list.Add(r.FormValue("description")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (t *ShoppingHandler) toggle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := t.list.ToggleDone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.todoView.View(w, item)
}

func (t *ShoppingHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := t.list.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *ShoppingHandler) reOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := t.list.ReOrder(r.Form["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	items, err := t.list.ShoppingItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.todoView.List(w, items)
}

func (t *ShoppingHandler) search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	results, err := t.list.Search(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.todoView.List(w, results)
}

func (t *ShoppingHandler) rename(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := t.list.Rename(id, r.Form["name"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.todoView.View(w, item)
}

func newStaticHandler() (http.Handler, error) {
	lol, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(lol)), nil
}
