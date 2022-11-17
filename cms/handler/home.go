package handler

import (
	// "log"
	"net/http"
)

type Auth struct {
	Auth	interface{}
}

func (s *Server) home(rw http.ResponseWriter, r *http.Request) {

	// session, err := h.sess.Get(r, sessionName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// auth := session.Values["authUserID"]
	// list := Auth{
	// 	Auth: auth,
	// }
	if err:= s.templates.ExecuteTemplate(rw, "home.html", nil); err != nil {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
	return
	}
}