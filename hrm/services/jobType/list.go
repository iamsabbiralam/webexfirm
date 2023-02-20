package jobType

import (
	"context"

	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) ListJobTypes(ctx context.Context, req *jobTypeG.ListJobTypesRequest) (*jobTypeG.ListJobTypesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.job-type.ListJobTypes")
	dbPrm := storage.JobTypes{
		SearchTerm: req.GetSearchTerm(),
		Offset:     req.GetOffset(),
		Limit:      req.GetLimit(),
		Status:     int(req.GetStatus()),
	}

	jobTypeList, err := h.store.ListJobTypes(ctx, dbPrm)
	if err != nil {
		errMsg := "failed to list job types"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	list := make([]*jobTypeG.JobTypes, len(jobTypeList))
	for key, val := range jobTypeList {
		list[key] = &jobTypeG.JobTypes{
			ID:        val.ID,
			Name:      val.Name,
			Status:    jobTypeG.Status(val.Status),
			Position:  int32(val.Position),
			CreatedAt: tspb.New(val.CreatedAt),
			UpdatedAt: tspb.New(val.UpdatedAt),
		}
	}

	var total int32
	if len(jobTypeList) > 0 {
		total = int32(jobTypeList[0].Count)
	}

	return &jobTypeG.ListJobTypesResponse{
		JobTypes: list,
		Total:    total,
	}, nil
}
