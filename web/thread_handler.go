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

type ThreadHandler struct {
	store    gothreadit.Store
	sessions *scs.SessionManager
}

// List
// Create
// Store
// Show
// Delete

func (th *ThreadHandler) List() http.HandlerFunc {
	type data struct {
		SessionData
		Threads []gothreadit.Thread
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := th.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{
			SessionData: GetSessionData(th.sessions, r.Context()),
			Threads:     tt,
		})
	}
}

func (th *ThreadHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			SessionData: GetSessionData(th.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (th *ThreadHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF   template.HTML
		Thread gothreadit.Thread
		Posts  []gothreadit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thead.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		t, err := th.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pp, err := th.store.PostsByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tmpl.Execute(w, data{
			GetSessionData(th.sessions, r.Context()),
			csrf.TemplateField(r),
			t,
			pp,
		})
	}
}

func (th *ThreadHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := CreateThreadForm{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Errors:      nil,
		}
		if !form.Validate() {
			th.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		err := th.store.CreateThread(&gothreadit.Thread{
			uuid.New(),
			form.Title,
			form.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		th.sessions.Put(r.Context(), "flash", "Your new thread has been created.")

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (th *ThreadHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = th.store.DeleteThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		th.sessions.Put(r.Context(), "flash", "The thread has been deleted")

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
