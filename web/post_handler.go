package server

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/rojerdu-dev/gothreadit"
	"html/template"
	"net/http"
)

type PostHandler struct {
	store    gothreadit.Store
	sessions *scs.SessionManager
}

// Create
// Store
// Show
// Vote

func (ph *PostHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData

		CSRF   template.HTML
		Thread gothreadit.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := ph.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			GetSessionData(ph.sessions, r.Context()),
			csrf.TemplateField(r),
			t,
		})
	}
}

func (ph *PostHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF     template.HTML
		Thread   gothreadit.Thread
		Post     gothreadit.Post
		Comments []gothreadit.Comment
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		threadIDStr := chi.URLParam(r, "threadID")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := ph.store.Post(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cc, err := ph.store.CommentsByPost(p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := ph.store.Thread(threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(ph.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			Thread:      t,
			Post:        p,
			Comments:    cc,
		})
	}
}

func (ph *PostHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := CreatePostForm{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
		}
		if !form.Validate() {
			ph.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := ph.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := &gothreadit.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    form.Title,
			Content:  form.Content,
		}

		err = ph.store.CreatePost(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ph.sessions.Put(r.Context(), "flash", "Your new post has been created")

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (ph *PostHandler) Vote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := ph.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			p.Votes++
		} else if dir == "down" {
			p.Votes--
		}

		err = ph.store.UpdatePost(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
