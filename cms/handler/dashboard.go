package handler

import (
	"net/http"
	"practice/webex/serviceutil/logging"
)

type DashboardTempData struct {
	FormErrors  map[string]string
	FormMessage map[string]string
}

func (s *Server) loadDashboardForm(w http.ResponseWriter, r *http.Request, data DashboardTempData) {
	if err := s.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getDashboardMethods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getDashboardMethods")
	data := DashboardTempData{}
	s.loadDashboardForm(w, r, data)
}
