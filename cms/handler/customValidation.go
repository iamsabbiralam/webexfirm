package handler

import (
	"context"
	"fmt"
	"net/http"

	jobTypeG "practice/webex/gunk/v1/jobType"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func validateJobTypesPosition(s *Server, r *http.Request, value int32, id string) validation.Rule {
	return validation.By(func(interface{}) error {
		dpts, _ := s.job.ListJobTypes(r.Context(), &jobTypeG.ListJobTypesRequest{})
		for _, item := range dpts.GetJobTypes() {
			if id != "" && item.GetID() == id && item.GetPosition() == value {
				return nil
			}

			if item.GetPosition() == value {
				return fmt.Errorf(" Position number is already exists")
			}
		}

		return nil
	})
}

func validateJobTypesName(s *Server, value string, id string) validation.Rule {
	return validation.By(func(interface{}) error {
		res, _ := s.job.GetJobTypesByTitle(context.Background(), &jobTypeG.GetJobTypesByTitleRequest{
			Title: value,
		})
		if id != "" && res != nil && res.GetID() == id && res.GetName() == value {
			return nil
		}

		if res == nil {
			return nil
		}

		return fmt.Errorf(" Job type already exists")
	})
}
