package jobType

import (
	"context"
	
	"practice/webex/hrm/storage"
)

type jobTypeStore interface {
	CreateJobType(context.Context, storage.JobTypes) (string, error)
	ListJobTypes(context.Context, storage.JobTypes) ([]storage.JobTypes, error)
	GetJobType(context.Context, string) (*storage.JobTypes, error)
	UpdateJobType(context.Context, storage.JobTypes) (string, error)
	DeleteJobType(context.Context, storage.JobTypes) error
}

type CoreSvc struct {
	store jobTypeStore
}

func NewCoreSvc(ccs jobTypeStore) *CoreSvc {
	return &CoreSvc{
		store: ccs,
	}
}
