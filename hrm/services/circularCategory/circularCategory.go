package circularCategory

import (
	"context"

	"practice/webex/hrm/storage"

	cc "practice/webex/gunk/v1/circularCategory"
)

type Handler struct {
	cc.UnimplementedCircularCategoryServiceServer
	ccst CircularCategoryStore
}

func New(cs CircularCategoryStore) *Handler {
	return &Handler{
		ccst: cs,
	}
}

type CircularCategoryStore interface {
	CreateCircularCategory(context.Context, storage.CircularCategory) (string, error)
	GetCircularCategory(context.Context, string) (*storage.CircularCategory, error)
	UpdateCircularCategory(context.Context, storage.CircularCategory) (*storage.CircularCategory, error)
	ListCircularCategory(context.Context, storage.CircularCategory) ([]storage.CircularCategory, error)
	DeleteCircularCategory(context.Context, storage.CircularCategory) error
}
