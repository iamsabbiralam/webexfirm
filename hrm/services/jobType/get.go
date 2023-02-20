package jobType

import (
	"context"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (h *Handler) GetJobTypes(ctx context.Context, req *jobTypeG.GetJobTypesRequest) (*jobTypeG.GetJobTypesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.job-type.GetJobTypes")
	res, err := h.store.GetJobType(ctx, req.GetID())
	if err != nil {
		errMsg := "no job type found"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &jobTypeG.GetJobTypesResponse{
		ID:       res.ID,
		Name:     res.Name,
		Status:   jobTypeG.Status(res.Status),
		Position: int32(res.Position),
	}, nil
}
