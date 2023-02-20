package jobType

import (
	"context"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (h *Handler) CreateJobTypes(ctx context.Context, req *jobTypeG.CreateJobTypesRequest) (*jobTypeG.CreateJobTypesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.Job-types.CreateJobTypes")
	dbPrm := storage.JobTypes{
		Name:     req.GetName(),
		Status:   int(req.GetStatus()),
		Position: int(req.GetPosition()),
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedAt: req.GetCreatedAt().AsTime(),
			CreatedBy: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
		},
	}
	id, err := h.store.CreateJobType(ctx, dbPrm)
	if err != nil {
		errMsg := "failed to create job types"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &jobTypeG.CreateJobTypesResponse{
		ID: id,
	}, nil
}
