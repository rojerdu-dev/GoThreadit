package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

func NewPostStore(db *sqlx.DB) *PostStore {
	return &PostStore{
		DB: db,
	}
}

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) Post(id uuid.UUID) (gothreadit.Post, error) {
	var p gothreadit.Post
	err := s.Get(&p, `SELECT * FROM posts WHERE id = $1`, id)
	if err != nil {
		return gothreadit.Post{}, fmt.Errorf("error getting post: %w", err)
	}
	return p, nil
}

func (s *PostStore) PostsByThread(threadID uuid.UUID) ([]gothreadit.Post, error) {
	var pp []gothreadit.Post
	err := s.Select(&pp, `SELECT * FROM posts WHERE thread_id = $1`, threadID)
	if err != nil {
		return []gothreadit.Post{}, fmt.Errorf("error getting posts: %w", err)
	}
	return pp, nil
}

func (s *PostStore) CreatePost(p *gothreadit.Post) error {
	err := s.Get(p, `INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
	)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}
	return nil
}

func (s *PostStore) UpdatePost(p *gothreadit.Post) error {
	err := s.Get(p, `UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5 RETURNING *`,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	_, err := s.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}