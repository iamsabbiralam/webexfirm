package handler

import (
	"log"
	"net/http"
	"practice/webex/serviceutil/logging"
)

type DashboardTempData struct {
	FormErrors  map[string]string
	FormMessage map[string]string
	GlobalURLs  map[string]string
}

func (s *Server) loadDashboardForm(w http.ResponseWriter, r *http.Request, data DashboardTempData) {
	template := s.lookupTemplate("dashboard.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusSeeOther)
		return
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, err.Error(), http.StatusSeeOther)
	}
}

func (s *Server) getDashboardMethods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getDashboardMethods")
	data := DashboardTempData{
		GlobalURLs: adminViewURLs(),
	}
	s.loadDashboardForm(w, r, data)
}
