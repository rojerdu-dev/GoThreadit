package postgres

import (
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
	panic("not implemented") // TODO: implement
}

func (s *ThreadStore) Threads() ([]gothreadit.Thread, error) {
	panic("not implemented") // TODO: implement
}

func (s *ThreadStore) CreateThread(t *gothreadit.Thread) error {
	panic("not implemented") // TODO: implement
}

func (s *ThreadStore) UpdateThread(t *gothreadit.Thread) error {
	panic("not implemented") // TODO: implement
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	panic("not implemented") // TODO: implement
}
