package server

import "encoding/gob"

func init() {
	gob.Register(CreatePostForm{})
	gob.Register(FormErrors{})
}

type FormErrors map[string]string

type CreatePostForm struct {
	Title   string
	Content string

	Errors FormErrors
}

func (f *CreatePostForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Title == "" {
		f.Errors["Title"] = "Please enter a title."
	}
	if f.Content == "" {
		f.Errors["Content"] = "Please enter a text."
	}

	return len(f.Errors) == 0
}
