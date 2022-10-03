package handler

import (
	"html/template"
	"net/http"
	"personal/webex/serviceutil/logging"
)

type DashboardTempData struct {
	CSRFField   template.HTML
	FormErrors  map[string]string
	FormMessage map[string]string
}

func (s *Server) loadDashboardForm(w http.ResponseWriter, r *http.Request, data DashboardTempData) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "getLoadDashboardHandler")
	template := s.templates.Lookup("dashboard.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, data); err != nil {
		errMsg := "error with template execution"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (s *Server) getDashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getDashboardHandler")
	data := DashboardTempData{}
	s.loadDashboardForm(w, r, data)
}
