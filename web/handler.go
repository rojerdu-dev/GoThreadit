package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rojerdu-dev/gothreadit"
	"html/template"
	"net/http"
)

func NewHandler(store gothreadit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	threads := ThreadHandler{store}
	posts := PostHandler{store}
	comments := CommentsHandler{store}

	h.Use(middleware.Logger)

	h.Get("/", h.Home())
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", threads.List())
		r.Get("/new", threads.Create())
		r.Post("/", threads.Store())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}", threads.Delete())
		r.Get("/{id}/new", posts.Create())
		r.Post("/{id}", posts.Store())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Get("/{threadID}/{postID}/vote", posts.Vote())
		r.Post("/{threadID}/{postID}", comments.Store())
	})
	h.Get("/comments/{id}/vote", comments.Vote())

	h.Get("/html", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("layout.html").ParseGlob("templates/includes/*.html"))
		t = template.Must(t.ParseFiles("templates/layout.html", "templates/childtemplate.html"))

		type params struct {
			Title   string
			Text    string
			Lines   []string
			Number1 int
			Number2 int
		}

		t.Execute(w, params{
			"Reddit Clone",
			"Welcome to Go ThreadIt Reddit Clone",
			[]string{
				"Line1",
				"Line2",
				"Line3",
			},
			49,
			53,
		})
	})

	return h
}

type Handler struct {
	*chi.Mux
	store gothreadit.Store
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		Posts []gothreadit.Post
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl.Execute(w, data{pp})
	}
}
