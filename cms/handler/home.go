package handler

import (
	"net/http"

	"personal/webex/serviceutil/logging"
)

func (s *Server) getHomeHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "getHomeHandler")
	template := s.lookupTemplate("home.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, template); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, "/error.html", http.StatusSeeOther)
	}
}
