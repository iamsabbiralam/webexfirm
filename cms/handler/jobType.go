package handler

import (
	"html/template"
	"net/http"
	"regexp"
	"time"

	"practice/webex/cms/paginator"
	jobTypeG "practice/webex/gunk/v1/jobType"
	"practice/webex/serviceutil/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-module/carbon"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type JobTypes struct {
	ID          string
	Title       string
	Status      jobTypeG.Status
	Position    int
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   time.Time
	DeletedBy   string
	CreatedTime string
	UpdatedTime string
}

type JobTypesTempData struct {
	CSRFField        template.HTML
	Form             JobTypes
	Data             []JobTypes
	IsReadCount      int32
	SearchTerm       string
	FormAction       string
	FormErrors       map[string]string
	FormMessage      map[string]string
	PaginationData   paginator.Paginator
	PresetPermission map[string]map[string]bool
	URLs             map[string]string
	GlobalURLs       map[string]string
	Status           []Status
}

func jobTypeURLs() map[string]string {
	return map[string]string{
		"create": createJobTypePath,
		"list":   jobTypesPath,
		"update": updateJobTypePath,
		"delete": deleteJobTypePath,
	}
}

func (job JobTypes) Validate(server *Server, r *http.Request, id string) error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&job,
		validation.Field(&job.Title, vre("The Title is required"), validation.Length(3, 100), validateJobTypesName(server, job.Title, id),
			validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")).Error("Job type should be character")),
		validation.Field(&job.Position, validation.When(job.ID != "", vre("The position is required"), validation.Min(1).Error("Value must be greater than or equal 1"), validateJobTypesPosition(server, r, int32(job.Position), id))),
		validation.Field(&job.Status, vre("The status is required")),
	)
}

func (s *Server) loadJobTypesTemplate(w http.ResponseWriter, r *http.Request, data JobTypesTempData) {
	log := logging.FromContext(r.Context()).WithField("method", "handler.job-type.loadJobTypesTemplate")
	tmpl := s.lookupTemplate("job-type.html")
	if tmpl == nil {
		log.Error("unable to load tmpl")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	if err := tmpl.Execute(w, data); err != nil {
		errMsg := "error with template execution"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) createJobTypeHandler(w http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).WithField("method", "handler.job-type.createJobTypeHandler")
	data := JobTypesTempData{
		CSRFField:  csrf.TemplateField(r),
		FormAction: createJobTypePath,
		URLs:       jobTypeURLs(),
		GlobalURLs: adminViewURLs(),
		Status:     GetStatus(jobTypeG.Status_name),
	}

	s.loadJobTypesTemplate(w, r, data)
}

func (s *Server) createPostJobTypeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.job-type.createPostJobTypeHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form JobTypes
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		errMsg := "Decoding form"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	if err := form.Validate(s, r, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := JobTypesTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
			FormAction: createJobTypePath,
			URLs:       jobTypeURLs(),
			GlobalURLs: adminViewURLs(),
			Status:     GetStatus(jobTypeG.Status_name),
		}

		s.loadJobTypesTemplate(w, r, data)
		return
	}

	jtList, err := s.job.ListJobTypes(ctx, &jobTypeG.ListJobTypesRequest{})
	if err != nil {
		errMsg := "failed to get job types"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	var pos int32
	if len(jtList.JobTypes) == 0 {
		pos = 1
	} else {
		pos = jtList.JobTypes[0].Position + 1
	}

	if _, err := s.job.CreateJobTypes(ctx, &jobTypeG.CreateJobTypesRequest{
		Name:      form.Title,
		Status:    form.Status,
		Position:  int32(pos),
		CreatedBy: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	}); err != nil {
		errMsg := "error with create job types"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	formMSG := map[string]string{"SuccessMessage": form.Title + " Added Successfully"}
	data := JobTypesTempData{
		CSRFField:   csrf.TemplateField(r),
		FormAction:  createJobTypePath,
		URLs:        jobTypeURLs(),
		GlobalURLs:  adminViewURLs(),
		FormMessage: formMSG,
		Status:      GetStatus(jobTypeG.Status_name),
	}
	s.loadJobTypesTemplate(w, r, data)
}

func (s *Server) getAllJobTypeHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "getAllJobTypeHandler")
	template := s.lookupTemplate("job-types.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	queryString := GetQueryStringData(r, []string{}, false)
	res, err := s.job.ListJobTypes(r.Context(), &jobTypeG.ListJobTypesRequest{
		SearchTerm: queryString.SearchTerm,
		Limit:      limitPerPage,
		Offset:     queryString.Offset,
	})
	if err != nil {
		log.Error("unable to get list: ", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	jtList := make([]JobTypes, 0, len(res.JobTypes))
	for _, item := range res.JobTypes {
		dptAppendData := JobTypes{
			ID:          item.ID,
			Title:       item.Name,
			Status:      item.Status,
			Position:    int(item.Position),
			CreatedTime: carbon.Parse(item.CreatedAt.AsTime().Format("2006-01-02 15:04:05")).DiffForHumans(),
			UpdatedTime: carbon.Parse(item.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")).DiffForHumans(),
		}
		jtList = append(jtList, dptAppendData)
	}

	formMSG := map[string]string{}
	if queryString.SearchTerm != "" && len(res.JobTypes) > 0 {
		formMSG = map[string]string{"FoundMessage": "Data Found"}
	} else if queryString.SearchTerm != "" && len(res.JobTypes) == 0 {
		formMSG = map[string]string{"NotFoundMessage": "Data Not Found"}
	}

	data := JobTypesTempData{
		Data:        jtList,
		SearchTerm:  queryString.SearchTerm,
		FormMessage: formMSG,
		URLs:        jobTypeURLs(),
		GlobalURLs:  adminViewURLs(),
		Status:      GetStatus(jobTypeG.Status_name),
	}
	if len(jtList) > 0 {
		data.PaginationData = paginator.NewPaginator(int32(queryString.CurrentPage), limitPerPage, res.Total, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) updateJobTypesHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "handler.job-type.updateJobTypesHandler")
	vars := mux.Vars(r)
	id := vars["id"]
	res, err := s.job.GetJobTypes(r.Context(), &jobTypeG.GetJobTypesRequest{
		ID: id,
	})
	if err != nil {
		log.Infof("error with get job types: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	jobTypeData := JobTypes{
		ID:       res.ID,
		Title:    res.Name,
		Status:   res.Status,
		Position: int(res.Position),
	}

	data := JobTypesTempData{
		CSRFField:  csrf.TemplateField(r),
		FormAction: DynamicUrlSwitch(updateJobTypePath, map[string]string{"id": id}),
		Form:       jobTypeData,
		URLs:       jobTypeURLs(),
		GlobalURLs: adminViewURLs(),
		Status:     GetStatus(jobTypeG.Status_name),
	}
	s.loadJobTypesTemplate(w, r, data)
}

func (s *Server) updatePostJobTypesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.job-type.updatePostJobTypesHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form JobTypes
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		return
	}

	if err := form.Validate(s, r, form.ID); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := JobTypesTempData{
			CSRFField:  csrf.TemplateField(r),
			Form:       form,
			FormErrors: vErrs,
			FormAction: updateJobTypePath,
			URLs:       jobTypeURLs(),
			GlobalURLs: adminViewURLs(),
			Status:     GetStatus(jobTypeG.Status_name),
		}
		s.loadJobTypesTemplate(w, r, data)
		return
	}

	if _, err := s.job.UpdateJobTypes(ctx, &jobTypeG.UpdateJobTypesRequest{
		ID:        form.ID,
		Name:      form.Title,
		Status:    form.Status,
		Position:  int32(form.Position),
		UpdatedBy: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	}); err != nil {
		log.Infof("error with update job types: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, jobTypesPath, http.StatusSeeOther)
}

func (s *Server) deleteJobTypeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.job-type.deleteJobTypeHandler")
	vars := mux.Vars(r)
	id := vars["id"]
	res, err := s.job.GetJobTypes(r.Context(), &jobTypeG.GetJobTypesRequest{
		ID: id,
	})
	if err != nil || res == nil {
		log.Infof("unable to get job types info: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	if _, err := s.job.DeleteJobTypes(ctx, &jobTypeG.DeleteJobTypesRequest{
		ID:        id,
		DeletedBy: "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	}); err != nil {
		log.Infof("error with delete job types: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, jobTypesPath, http.StatusSeeOther)
}
