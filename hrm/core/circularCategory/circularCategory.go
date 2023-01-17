package circularCategory

import (
	"context"

	"practice/webex/hrm/storage"
)

type circularCategoryStore interface {
	CreateCircularCategory(context.Context, storage.CircularCategory) (string, error)
	GetCircularCategory(context.Context, string) (*storage.CircularCategory, error)
	UpdateCircularCategory(context.Context, storage.CircularCategory) (*storage.CircularCategory, error)
	ListCircularCategory(context.Context, storage.CircularCategory) ([]storage.CircularCategory, error)
	DeleteCircularCategory(context.Context, storage.CircularCategory) error
}

type CoreSvc struct {
	store circularCategoryStore
}

func NewCoreSvc(ccs circularCategoryStore) *CoreSvc {
	return &CoreSvc{
		store: ccs,
	}
}
