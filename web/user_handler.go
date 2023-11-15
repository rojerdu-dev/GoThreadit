package server

import (
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/rojerdu-dev/gothreadit"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserHandler struct {
	store    gothreadit.Store
	sessions *scs.SessionManager
}

func (uh *UserHandler) Register() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/user_register.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			GetSessionData(uh.sessions, r.Context()),
			csrf.TemplateField(r),
		})
	}
}

func (uh *UserHandler) RegisterSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := RegisterForm{
			Username:      r.FormValue("username"),
			Password:      r.FormValue("password"),
			UsernameTaken: false,
		}
		_, err := uh.store.UsersByUsername(form.Username)
		if err == nil {
			form.UsernameTaken = true
		}
		if !form.Validate() {
			uh.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = uh.store.CreateUser(&gothreadit.User{
			ID:       uuid.UUID{},
			Username: form.Username,
			Password: string(password),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uh.sessions.Put(r.Context(), "flash", "Your registration was successful. Please log in.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (uh *UserHandler) Login() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/user_login.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			GetSessionData(uh.sessions, r.Context()),
			csrf.TemplateField(r),
		})
	}
}

func (uh *UserHandler) LoginSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := LoginForm{
			Username:             r.FormValue("username"),
			Password:             r.FormValue("password"),
			IncorrectCredentials: false,
		}

		user, err := uh.store.UsersByUsername(form.Username)
		if err != nil {
			form.IncorrectCredentials = true
		} else {
			compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
			form.IncorrectCredentials = compareErr != nil
		}
		if !form.Validate() {
			uh.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		uh.sessions.Put(r.Context(), "user_id", user.ID)
		uh.sessions.Put(r.Context(), "flash", "Login successful.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (uh *UserHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uh.sessions.Remove(r.Context(), "user_id")
		uh.sessions.Put(r.Context(), "flash", "You have been logged out successfully.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
