package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

func NewThreadStore(db *sqlx.DB) *ThreadStore {
	return &ThreadStore{
		DB: db,
	}
}

type ThreadStore struct {
	*sqlx.DB
}

func (s *ThreadStore) Thread(id uuid.UUID) (gothreadit.Thread, error) {
	var t gothreadit.Thread
	err := s.Get(&t, `SELECT * FROm threads WHERE id = $1`, id)
	if err != nil {
		return gothreadit.Thread{}, fmt.Errorf("error getting thread: %w", err)
	}
	return t, nil
}

func (s *ThreadStore) Threads() ([]gothreadit.Thread, error) {
	var tt []gothreadit.Thread
	err := s.Select(&tt, `SELECT * FROM threads`)
	if err != nil {
		return []gothreadit.Thread{}, fmt.Errorf("error getting threads: %w", err)
	}
	return tt, nil
}

func (s *ThreadStore) CreateThread(t *gothreadit.Thread) error {
	err := s.Get(t, `INSERT INTO threads VALUES ($1, $2, $3) RETURNING *`,
		t.Title,
		t.Description,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("error creating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) UpdateThread(t *gothreadit.Thread) error {
	err := s.Get(t, `UPDATE threads SET title = $1, description = $2 WHERE id = $3 RETURNING *`,
		t.ID,
		t.Title,
		t.Description,
	)
	if err != nil {
		return fmt.Errorf("error updating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	_, err := s.Exec(`DELETE FROM threads WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting thread: %w", err)
	}
	return nil
}
