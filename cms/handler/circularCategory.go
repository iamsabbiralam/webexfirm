package handler

import (
	"html/template"
	"net/http"
	"time"

	"practice/webex/cms/paginator"
	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/serviceutil/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-module/carbon"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CircularCategory struct {
	ID          string
	Name        string
	Description string
	Status      cc.Status
	Position    int32
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	CreatedTime string
	UpdatedTime string
}

type CircularCategoryTempData struct {
	CSRFField      template.HTML
	Form           CircularCategory
	FormAction     string
	Data           []CircularCategory
	SearchTerm     string
	Status         []Status
	FormMessage    map[string]string
	FormErrors     map[string]string
	PaginationData paginator.Paginator
	URLs           map[string]string
	GlobalURLs     map[string]string
}

func circularCategoryURLs() map[string]string {
	return map[string]string{
		"create":       createCircularCategoryPath,
		"list":         circularCategoriesPath,
		"update":       updateCircularCategoryPath,
		"delete":       deleteCircularCategoryPath,
		"updateStatus": updateCircularCategoryStatusPath,
	}
}

func (s CircularCategory) ValidateCircularCategory() error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, vre("The name is required")),
		validation.Field(&s.Description, vre("The description is required")),
		validation.Field(&s.Status, vre("The status is required")),
		validation.Field(&s.Position, validation.When(s.ID != "", vre("The position is required"), validation.Min(1).Error("Put valid position number"))),
	)
}

func (s *Server) loadCircularCategoryTemplate(w http.ResponseWriter, r *http.Request, data CircularCategoryTempData) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.circular-category.loadCircularCategoryTemplate")
	tmpl := s.lookupTemplate("circular-category.html")
	if tmpl == nil {
		errMsg := "unable to create circular category template"
		log.Errorf(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) createCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "handler.circular-category.createCircularCategoryHandler")
	data := CircularCategoryTempData{
		CSRFField:  csrf.TemplateField(r),
		FormAction: createCircularCategoryPath,
		URLs:       circularCategoryURLs(),
		GlobalURLs: adminViewURLs(),
		Status:     GetStatus(cc.Status_name),
	}

	s.loadCircularCategoryTemplate(w, r, data)
}

func (s *Server) postCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.circular-category.postCircularCategoryHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, ErrorPath, http.StatusBadRequest)
		return
	}

	var form CircularCategory
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	errMessage := form.ValidateCircularCategory()
	data := CircularCategoryTempData{}
	if errMessage != nil {
		s.validateMsg(w, r, errMessage, data, form)
		return
	}

	// get latest position
	circularCategories, err := s.cc.ListCircularCategory(ctx, &cc.ListCircularCategoryRequest{})
	if err != nil {
		errMsg := "unable to get latest position"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	var pos int64
	if len(circularCategories.CircularCategory) == 0 {
		pos = 1
	} else {
		pos = circularCategories.CircularCategory[0].Position + 1
	}

	_, err = s.cc.CreateCircularCategory(ctx, &cc.CreateCircularCategoryRequest{
		Name:        form.Name,
		Description: form.Description,
		Status:      form.Status,
		Position:    int64(pos),
		CreatedBy:   "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	})
	if err != nil {
		errMsg := "error with create circular category"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	formMSG := map[string]string{"SuccessMessage": form.Name + " Added Successfully"}
	data = CircularCategoryTempData{
		CSRFField:   csrf.TemplateField(r),
		FormMessage: formMSG,
		FormAction:  createCircularCategoryPath,
		URLs:        circularCategoryURLs(),
		GlobalURLs:  adminViewURLs(),
		Status:      GetStatus(cc.Status_name),
	}
	s.loadCircularCategoryTemplate(w, r, data)
}

func (s *Server) listCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.circular-category.listCircularCategoryHandler")
	template := s.lookupTemplate("circular-categories.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	queryString := GetQueryStringData(r, []string{}, false)
	res, err := s.cc.ListCircularCategory(ctx, &cc.ListCircularCategoryRequest{
		SearchTerm: queryString.SearchTerm,
		Offset:     queryString.Offset,
		Limit:      limitPerPage,
	})
	if err != nil {
		errMsg := "unable to get list"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	circularCategories := []CircularCategory{}
	var ccTotal int32
	if res != nil {
		ccTotal = res.Total
		for _, item := range res.GetCircularCategory() {
			getAllAppendData := CircularCategory{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Position:    int32(item.Position),
				Status:      item.Status,
				CreatedTime: carbon.Parse(item.CreatedAt.AsTime().Format("2006-01-02 15:04:05")).DiffForHumans(),
				UpdatedTime: carbon.Parse(item.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")).DiffForHumans(),
			}
			circularCategories = append(circularCategories, getAllAppendData)
		}
	}

	formMessage := map[string]string{}
	// search message conditions
	if queryString.SearchTerm != "" && res != nil && len(res.GetCircularCategory()) > 0 {
		formMessage = map[string]string{"FoundMessage": "Data Found"}
	} else if queryString.SearchTerm != "" && res != nil && len(res.GetCircularCategory()) == 0 {
		formMessage = map[string]string{"NotFoundMessage": "Data Not Found"}
	}

	data := CircularCategoryTempData{
		CSRFField:   csrf.TemplateField(r),
		Data:        circularCategories,
		SearchTerm:  queryString.SearchTerm,
		FormMessage: formMessage,
		URLs:        circularCategoryURLs(),
		GlobalURLs:  adminViewURLs(),
		Status:      GetStatus(cc.Status_name),
	}
	paginatorData := paginator.NewPaginator(int32(queryString.CurrentPage), limitPerPage, ccTotal, r)
	if len(circularCategories) > 0 {
		data.PaginationData = paginatorData
	}

	if err := template.Execute(w, data); err != nil {
		log.Errorf("error with template execution: %+v", err)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) getCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "handler.circular-category.getCircularCategoryHandler")
	vars := mux.Vars(r)
	id := vars["id"]
	res, err := s.cc.GetCircularCategory(r.Context(), &cc.GetCircularCategoryRequest{
		ID: id,
	})
	if err != nil {
		errMsg := "error with get circular category by id"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	circularC := CircularCategory{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Position:    int32(res.Position),
		Status:      res.Status,
	}

	data := CircularCategoryTempData{
		CSRFField:  csrf.TemplateField(r),
		Form:       circularC,
		FormAction: DynamicUrlSwitch(updateCircularCategoryPath, map[string]string{"id": id}),
		URLs:       circularCategoryURLs(),
		GlobalURLs: adminViewURLs(),
		Status:     GetStatus(cc.Status_name),
	}

	s.loadCircularCategoryTemplate(w, r, data)
}

func (s *Server) updateCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.circular-category.updateCircularCategoryHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form CircularCategory
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		errMsg := "decoding form"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	errMessage := form.ValidateCircularCategory()
	data := CircularCategoryTempData{}
	if errMessage != nil {
		s.validateMsg(w, r, errMessage, data, form)
		return
	}

	_, err := s.cc.UpdateCircularCategory(ctx, &cc.UpdateCircularCategoryRequest{
		ID:          form.ID,
		Name:        form.Name,
		Description: form.Description,
		Position:    int64(form.Position),
		Status:      form.Status,
		UpdatedBy:   "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	})
	if err != nil {
		errMsg := "error with update circular category"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, circularCategoriesPath, http.StatusSeeOther)
}

func (s *Server) updateCircularCategoryStatusHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "handler.circular-category.updateCircularCategoryStatusHandler")
	vars := mux.Vars(r)
	id := vars["id"]
	res, err := s.cc.GetCircularCategory(r.Context(), &cc.GetCircularCategoryRequest{
		ID: id,
	})
	if err != nil {
		errMsg := "Failed to get value"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	status := cc.Status_Active
	if res.Status == cc.Status_Active {
		status = cc.Status_Inactive
	}

	_, err = s.cc.UpdateCircularCategory(r.Context(), &cc.UpdateCircularCategoryRequest{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      status,
		Position:    res.Position,
		UpdatedAt:   timestamppb.New(res.UpdatedAt.AsTime()),
		UpdatedBy:   "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	})

	if err != nil {
		errMsg := "error with Update circular category status"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, circularCategoriesPath, http.StatusSeeOther)
}

func (s *Server) deleteCircularCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "handler.circular-category.deleteCircularCategoryHandler")
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := s.cc.DeleteCircularCategory(r.Context(), &cc.DeleteCircularCategoryRequest{
		ID: id,
	})
	if err != nil {
		errMsg := "Failed to get value"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, circularCategoriesPath, http.StatusSeeOther)
}

func (s *Server) validateMsg(w http.ResponseWriter, r *http.Request, err error, data CircularCategoryTempData, form CircularCategory) error {
	vErrs := map[string]string{}
	if e, ok := err.(validation.Errors); ok {
		if len(e) > 0 {
			for key, value := range e {
				vErrs[key] = value.Error()
			}
		}
	}

	if form.ID != "" {
		data.FormAction = createCircularCategoryPath
	}

	data = CircularCategoryTempData{
		CSRFField:  csrf.TemplateField(r),
		Form:       form,
		FormErrors: vErrs,
		FormAction: DynamicUrlSwitch(updateCircularCategoryPath, map[string]string{"id": form.ID}),
		URLs:       circularCategoryURLs(),
		GlobalURLs: adminViewURLs(),
		Status:     GetStatus(cc.Status_name),
	}

	s.loadCircularCategoryTemplate(w, r, data)
	return nil
}
