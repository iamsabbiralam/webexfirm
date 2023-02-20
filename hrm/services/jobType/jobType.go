package jobType

import (
	"context"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/hrm/storage"
)

type Handler struct {
	jobTypeG.UnimplementedJobTypesServiceServer
	store JobTypesStore
}

type JobTypesStore interface {
	CreateJobType(context.Context, storage.JobTypes) (string, error)
	ListJobTypes(context.Context, storage.JobTypes) ([]storage.JobTypes, error)
	GetJobType(context.Context, string) (*storage.JobTypes, error)
	UpdateJobType(context.Context, storage.JobTypes) (string, error)
	DeleteJobType(context.Context, storage.JobTypes) error
}

func New(cs JobTypesStore) *Handler {
	return &Handler{
		store: cs,
	}
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Handler) JobTypes(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreateJobType": {res: "HRM:JobType", act: "create"},
		"ListJobType":   {res: "HRM:JobType", act: "view"},
		"DeleteJobType": {res: "HRM:JobType", act: "delete"},
		"UpdateJobType": {res: "HRM:JobType", act: "update"},
		"GetJobType":    {res: "HRM:JobType", act: "get"},
	}

	return p[mthd].res, p[mthd].act, p[mthd].pub
}
