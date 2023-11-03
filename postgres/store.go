package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rojerdu-dev/gothreadit"
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

	return &Store{
		NewThreadStore(db),
		NewPostStore(db),
		NewCommentStore(db),
	}, nil
}

type Store struct {
	gothreadit.ThreadStore
	gothreadit.PostStore
	gothreadit.CommentStore
}
