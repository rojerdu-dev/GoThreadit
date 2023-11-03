package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

func NewCommentStore(db *sqlx.DB) *CommentStore {
	return &CommentStore{
		DB: db,
	}
}

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) CommentsByPost() ([]gothreadit.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *CommentStore) Comment(id uuid.UUID) (gothreadit.Comment, error) {
	var c gothreadit.Comment
	err := s.Get(&c, `SELECT * FROM comments WHERE id = $1`, id)
	if err != nil {
		return gothreadit.Comment{}, fmt.Errorf("error getting comment: %w", err)
	}
	return c, nil
}

func (s *CommentStore) CommentByPost(postID uuid.UUID) ([]gothreadit.Comment, error) {
	var cc []gothreadit.Comment
	err := s.Select(&cc, `SELECT * FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		return []gothreadit.Comment{}, fmt.Errorf("error getting comments: %w", err)
	}
	return cc, nil
}

func (s *CommentStore) CreateComment(c *gothreadit.Post) error {
	err := s.Get(c, `INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *`,
		c.ID,
		c.PostID,
		c.Content,
		c.Votes,
	)
	if err != nil {
		return fmt.Errorf("error creating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(c *gothreadit.Post) error {
	err := s.Get(c, `UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETRUNING *`,
		c.PostID,
		c.Content,
		c.Votes,
		c.ID)
	if err != nil {
		fmt.Printf("error updating comments: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	_, err := s.Exec(`DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleing comment: %w", err)
	}
	return nil
}
