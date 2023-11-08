package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rojerdu-dev/gothreadit"
	"net/http"
)

type CommentsHandler struct {
	store gothreadit.Store
}

// Store
// Vote

func (ch *CommentsHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := r.FormValue("content")
		idStr := chi.URLParam(r, "PostID")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c := &gothreadit.Comment{
			ID:      uuid.New(),
			PostID:  id,
			Content: content,
		}
		err = ch.store.CreateComment(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (ch *CommentsHandler) Vote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c, err := ch.store.Comment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			c.Votes++
		} else if dir == "down" {
			c.Votes--
		}

		err = ch.store.UpdateComment(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
