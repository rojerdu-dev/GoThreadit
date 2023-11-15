package server

import (
	"context"
	"database/sql"
	"encoding/gob"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/rojerdu-dev/gothreadit"
)

func init() {
	gob.Register(uuid.UUID{})
}

func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	sessions := scs.New()
	sessions.Store = postgresstore.New(db)

	return sessions, nil
}

type SessionData struct {
	FlashMessage string
	Form         interface{}
	User         gothreadit.User
	LoggedIn     bool
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	data.FlashMessage = session.PopString(ctx, "flash")
	data.User, data.LoggedIn = ctx.Value("user").(gothreadit.User)

	data.Form = session.Pop(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	return data
}
