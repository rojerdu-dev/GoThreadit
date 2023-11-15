package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	//return &Store{
	//	&ThreadStore{db},
	//	&PostStore{db},
	//	&CommentStore{db},
	//}, nil
	return &Store{
		ThreadStore:  &ThreadStore{db},
		PostStore:    &PostStore{db},
		CommentStore: &CommentStore{db},
		UserStore:    &UserStore{db},
	}, nil

}

type Store struct {
	*ThreadStore
	*PostStore
	*CommentStore
	*UserStore
}
