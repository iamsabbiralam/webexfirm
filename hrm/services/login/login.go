package login

import (
	"context"

	login "personal/webex/gunk/v1/login"
	"personal/webex/hrm/storage"
)

type Handler struct {
	login.UnimplementedLoginServiceServer
	ls LoginStore
}

type LoginStore interface {
	Login(context.Context, string) (storage.SignUP, error)
}

func New(ls LoginStore) *Handler {
	return &Handler{
		ls: ls,
	}
}
