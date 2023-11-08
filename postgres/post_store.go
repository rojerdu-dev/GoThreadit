package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

type PostStore struct {
	*sqlx.DB
}

func (ps *PostStore) Post(id uuid.UUID) (gothreadit.Post, error) {
	var p gothreadit.Post
	err := ps.Get(&p, `SELECT * FROM posts WHERE id = $1`, id)
	if err != nil {
		return gothreadit.Post{}, fmt.Errorf("error getting post: %w", err)
	}
	return p, nil
}

func (ps *PostStore) PostsByThread(threadID uuid.UUID) ([]gothreadit.Post, error) {
	var pp []gothreadit.Post
	var query = `
		SELECT 
		    posts.*,
			COUNT(comments.*) AS comments_count
		FROM posts 
		JOIN comments ON comments.post_id = posts.id
		WHERE thread_id = $1
		GROUP BY posts.id
		ORDER BY votes DESC`

	err := ps.Select(&pp, query, threadID)
	if err != nil {
		return []gothreadit.Post{}, fmt.Errorf("error getting posts: %w", err)
	}
	return pp, nil
}

func (ps *PostStore) Posts() ([]gothreadit.Post, error) {
	var pp []gothreadit.Post
	var query = `
		SELECT 
		    posts.*,
			COUNT(comments.*) AS comments_count,
			threads.title AS thread_title
		FROM posts 
		JOIN comments ON comments.post_id = posts.id
		JOIN threads ON threads.id = posts.thread_id
		GROUP BY posts.id, threads.title
		ORDER BY votes DESC`

	err := ps.Select(&pp, query)
	if err != nil {
		return []gothreadit.Post{}, fmt.Errorf("error getting posts: %w", err)
	}
	return pp, nil
}

func (ps *PostStore) CreatePost(p *gothreadit.Post) error {
	err := ps.Get(p, `INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *`,
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

func (ps *PostStore) UpdatePost(p *gothreadit.Post) error {
	err := ps.Get(p, `UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5 RETURNING *`,
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

func (ps *PostStore) DeletePost(id uuid.UUID) error {
	_, err := ps.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}
