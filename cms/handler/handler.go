package handler

import (
	"errors"
	"io/fs"
	"log"
	"net/http"

	// "text/template"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	user "practice/webex/gunk/v1/user"
	"practice/webex/serviceutil/mw"
)

const sessionName = "webex-session"

func Handler(
	decoder *schema.Decoder,
	config *viper.Viper,
	logger *logrus.Entry,
	sess *sessions.CookieStore,
	hrmConn *grpc.ClientConn,
	assets fs.FS,
) (*mux.Router, error) {
	s := &Server{
		decoder: decoder,
		sess:    sess,
		assets:  assets,
		assetFS: hashfs.NewFS(assets),
		user: struct {
			user.UserServiceClient
		}{
			UserServiceClient: user.NewUserServiceClient(hrmConn),
		},
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	csrfSecure := config.GetBool("csrf.secure")
	csrfSecret := config.GetString("csrf.secret")
	if csrfSecret == "" {
		return nil, errors.New("CSRF secret must not be empty")
	}

	r := mux.NewRouter()
	r.HandleFunc(homeURL, s.home)
	mw.ChainHTTPMiddleware(r, logger,
		mw.CSRF([]byte(csrfSecret), csrf.Secure(csrfSecure), csrf.Path("/")),
	)

	r.HandleFunc(registrationURL, s.signUpMethod).Methods("GET")
	r.HandleFunc(registrationURL, s.postSignUpMethod).Methods("POST")
	r.HandleFunc(loginURL, s.getLoginHandler).Methods("GET")
	r.HandleFunc(loginURL, s.postLoginHandler).Methods("POST")
	r.HandleFunc(logoutPath, s.logoutHandler).Methods("GET").Name("logout")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", cacheStaticFiles(http.FileServer(http.FS(s.assetFS)))))

	l := r.NewRoute().Subrouter()
	l.Use(s.loginMiddleware)

	m := r.NewRoute().Subrouter()
	m.Use(s.authMiddleware)
	m.HandleFunc(dashboardPath, s.getDashboardMethods).Methods("GET").Name("dashboard")
	m.HandleFunc(getAllUsersPath, s.getAllUsersHandler).Methods("GET").Name("userList")

	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := s.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, "invalid URL", http.StatusInternalServerError)
			return
		}
	})

	return r, nil
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
