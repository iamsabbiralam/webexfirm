package signUp

import (
	"context"

	signup "personal/webex/gunk/v1/signUp"
	"personal/webex/hrm/storage"
)

type Handler struct {
	signup.UnimplementedSignUpServiceServer
	su SignUPStore
}

type SignUPStore interface {
	Register(context.Context, storage.SignUP) (string, error)
}

func New(ss SignUPStore) *Handler {
	return &Handler{
		su: ss,
	}
}
