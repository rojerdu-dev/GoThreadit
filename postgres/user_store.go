package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rojerdu-dev/gothreadit"
)

type UserStore struct {
	*sqlx.DB
}

func (us *UserStore) User(id uuid.UUID) (gothreadit.User, error) {
	var t gothreadit.User
	err := us.Get(&t, `SELECT * FROm users WHERE id = $1`, id)
	if err != nil {
		return gothreadit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return t, nil
}

func (us *UserStore) UsersByUsername(username string) (gothreadit.User, error) {
	var u gothreadit.User
	err := us.Select(&u, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		return gothreadit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (us *UserStore) Users() ([]gothreadit.User, error) {
	var uu []gothreadit.User
	err := us.Select(&uu, `SELECT * FROM users`)
	if err != nil {
		return []gothreadit.User{}, fmt.Errorf("error getting users: %w", err)
	}
	return uu, nil
}

func (us *UserStore) CreateUser(u *gothreadit.User) error {
	err := us.Get(u, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		u.ID,
		u.Username,
		u.Password,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (us *UserStore) UpdateUser(u *gothreadit.User) error {
	err := us.Get(u, `UPDATE users SET username = $1, username = $2 WHERE id = $3 RETURNING *`,
		u.Username,
		u.Password,
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating thread: %w", err)
	}
	return nil
}

func (us *UserStore) DeleteUser(id uuid.UUID) error {
	_, err := us.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
