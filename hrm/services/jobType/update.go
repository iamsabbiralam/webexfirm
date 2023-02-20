package jobType

import (
	"context"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateJobTypes(ctx context.Context, req *jobTypeG.UpdateJobTypesRequest) (*jobTypeG.UpdateJobTypesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.job-type.UpdateJobTypes")
	dbPrm := storage.JobTypes{
		ID:       req.GetID(),
		Name:     req.GetName(),
		Status:   int(req.GetStatus()),
		Position: int(req.GetPosition()),
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
		},
	}
	id, err := h.store.UpdateJobType(ctx, dbPrm)
	if err != nil {
		errMsg := "no job type found"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &jobTypeG.UpdateJobTypesResponse{
		ID: id,
	}, nil
}
