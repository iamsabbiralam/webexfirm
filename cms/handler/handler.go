package handler

import (
	"log"
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
	
	l := r.NewRoute().Subrouter()
	l.HandleFunc(registrationURL, s.signUpMethod).Methods("GET")
	l.HandleFunc(registrationURL, s.postSignUpMethod).Methods("POST")
	l.HandleFunc(loginURL, s.getLoginHandler).Methods("GET")
	l.HandleFunc(loginURL, s.postLoginHandler).Methods("POST")
	l.Use(s.loginMiddleware)

	m := r.NewRoute().Subrouter()
	m.Use(s.authMiddleware)
	m.HandleFunc(dashboardPath, s.getDashboardMethods).Methods("GET")

	m.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))
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
		"cms/assets/templates/register/create-user.html",
		"cms/assets/templates/register/login.html",
		"cms/assets/templates/user/dashboard.html",
	))

	return nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sess.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserID"]
		if authUserID != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		}
		
	})
}

func (s *Server) loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sess.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserID"]
		if authUserID != nil {
			http.Redirect(w, r, homeURL, http.StatusTemporaryRedirect)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
