package handler

import (
	"log"
	"net/http"
	"practice/webex/serviceutil/logging"
)

type Auth struct {
	Auth	interface{}
}

type HomeTempData struct {
	Auth	Auth
	GlobalURLs	map[string]string
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	// session, err := h.sess.Get(r, sessionName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// auth := session.Values["authUserID"]
	// list := Auth{
	// 	Auth: auth,
	// }
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getLoginHandler")
	data := HomeTempData{
		GlobalURLs: adminViewURLs(),
	}
	s.loadHomePage(w, r, data)
}

func (s *Server) loadHomePage(w http.ResponseWriter, r *http.Request, data HomeTempData) {
	template := s.lookupTemplate("home.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusSeeOther)
		return
	}
	
	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, "", http.StatusSeeOther)
	}
}
