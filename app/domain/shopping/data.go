package shopping

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS shopping_items (
            id TEXT PRIMARY KEY,
            description TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            complete BOOLEAN NOT NULL DEFAULT FALSE
        );
    `)
	return err
}

func (s *SQLiteStore) Add(item ShoppingItem) error {
	_, err := s.db.Exec(
		"INSERT INTO shopping_items (id, description, created_at, complete) VALUES (?, ?, ?, ?)",
		item.ID.String(), item.Description, item.CreatedAt, item.Complete,
	)
	return err
}

func (s *SQLiteStore) GetAll() ([]ShoppingItem, error) {
	rows, err := s.db.Query("SELECT id, description, created_at, complete FROM shopping_items ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ShoppingItem
	for rows.Next() {
		var item ShoppingItem
		var id string
		if err := rows.Scan(&id, &item.Description, &item.CreatedAt, &item.Complete); err != nil {
			return nil, err
		}
		item.ID, _ = uuid.Parse(id)
		items = append(items, item)
	}
	return items, nil
}

func (s *SQLiteStore) Update(item ShoppingItem) error {
	_, err := s.db.Exec(
		"UPDATE shopping_items SET description = ?, complete = ? WHERE id = ?",
		item.Description, item.Complete, item.ID.String(),
	)
	return err
}

func (s *SQLiteStore) Delete(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM shopping_items WHERE id = ?", id.String())
	return err
}

func (s *SQLiteStore) ReOrder(ids []uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, id := range ids {
		_, err := tx.Exec(
			"UPDATE shopping_items SET created_at = ? WHERE id = ?",
			time.Now().Add(time.Duration(i)*time.Millisecond),
			id.String(),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
