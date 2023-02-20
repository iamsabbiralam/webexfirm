package jobType

import (
	"context"
	"database/sql"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (h *Handler) DeleteJobTypes(ctx context.Context, req *jobTypeG.DeleteJobTypesRequest) (*jobTypeG.DeleteJobTypesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.job-type.DeleteJobTypes")
	dbPrm := storage.JobTypes{
		ID: req.GetID(),
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{
				String: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
				Valid:  true,
			},
		},
	}
	err := h.store.DeleteJobType(ctx, dbPrm)
	if err != nil {
		errMsg := "failed to delete by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &jobTypeG.DeleteJobTypesResponse{}, nil
}
