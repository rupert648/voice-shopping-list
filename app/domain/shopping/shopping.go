package shopping

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ShoppingItem struct {
	ID          uuid.UUID
	Description string
	CreatedAt   time.Time
	Complete    bool
}

type List struct {
	store *SQLiteStore
}

func NewList(dbPath string) (*List, error) {
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}
	return &List{store: store}, nil
}

func (s *List) Add(description string) error {
	item := ShoppingItem{
		ID:          uuid.New(),
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.store.Add(item)
}

func (s *List) Rename(id uuid.UUID, name string) (ShoppingItem, error) {
	items, err := s.store.GetAll()
	if err != nil {
		return ShoppingItem{}, err
	}

	for _, item := range items {
		if item.ID == id {
			item.Description = name
			err := s.store.Update(item)
			return item, err
		}
	}
	return ShoppingItem{}, fmt.Errorf("item not found")
}

func (s *List) ShoppingItems() ([]ShoppingItem, error) {
	return s.store.GetAll()
}

func (s *List) ToggleDone(id uuid.UUID) (ShoppingItem, error) {
	items, err := s.store.GetAll()
	if err != nil {
		return ShoppingItem{}, err
	}

	for _, item := range items {
		if item.ID == id {
			item.Complete = !item.Complete
			err := s.store.Update(item)
			return item, err
		}
	}
	return ShoppingItem{}, fmt.Errorf("item not found")
}

func (s *List) Delete(id uuid.UUID) error {
	return s.store.Delete(id)
}

func (s *List) ReOrder(ids []string) error {
	var uuids []uuid.UUID
	for _, id := range ids {
		uuids = append(uuids, uuid.MustParse(id))
	}
	return s.store.ReOrder(uuids)
}

func (s *List) Search(search string) ([]ShoppingItem, error) {
	items, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}

	search = strings.ToLower(search)
	var results []ShoppingItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Description), search) {
			results = append(results, item)
		}
	}
	return results, nil
}
