package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

type CommentStore struct {
	*sqlx.DB
}

func (cs *CommentStore) Comment(id uuid.UUID) (gothreadit.Comment, error) {
	var c gothreadit.Comment
	err := cs.Get(&c, `SELECT * FROM comments WHERE id = $1`, id)
	if err != nil {
		return gothreadit.Comment{}, fmt.Errorf("error getting comment: %w", err)
	}
	return c, nil
}

func (cs *CommentStore) CommentsByPost(postID uuid.UUID) ([]gothreadit.Comment, error) {
	var cc []gothreadit.Comment
	err := cs.Select(&cc, `SELECT * FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		return []gothreadit.Comment{}, fmt.Errorf("error getting comments: %w", err)
	}
	return cc, nil
}

func (cs *CommentStore) CreateComment(c *gothreadit.Comment) error {
	err := cs.Get(c, `INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *`,
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

func (s *CommentStore) UpdateComment(c *gothreadit.Comment) error {
	err := s.Get(c, `UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *`,
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
