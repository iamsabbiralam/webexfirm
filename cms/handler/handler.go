package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	user "practice/webex/gunk/v1/user"
)

const sessionName = "webex-session"

func Handler(
	decoder *schema.Decoder,
	config *viper.Viper,
	sess *sessions.CookieStore,
	hrmConn *grpc.ClientConn,
) (*mux.Router, error) {
	s := &Server{
		decoder: decoder,
		sess:    sess,
		user: struct {
			user.UserServiceClient
		}{
			UserServiceClient: user.NewUserServiceClient(hrmConn),
		},
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.HandleFunc(homeURL, s.home)
	r.HandleFunc(registrationURL, s.signUpMethod).Methods("GET")
	r.HandleFunc(registrationURL, s.postSignUpMethod).Methods("POST")
	r.HandleFunc(loginURL, s.getLoginHandler).Methods("GET")
	r.HandleFunc(loginURL, s.postLoginHandler).Methods("POST")
	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))
	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := s.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, "invalid URL", http.StatusInternalServerError)
			return
		}
	})

	return r, nil
}

func (s *Server) parseTemplates() error {
	s.templates = template.Must(template.ParseFiles(
		"cms/assets/templates/base/home.html",
		"cms/assets/templates/user/create-user.html",
		"cms/assets/templates/user/login.html",
	))

	return nil
}
