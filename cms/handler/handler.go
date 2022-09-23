package handler

import (
	"errors"
	"io/fs"
	"net/http"

	logING "personal/webex/gunk/v1/login"
	signUpG "personal/webex/gunk/v1/signUp"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"personal/webex/serviceutil/mw"
)

const sessionName = "webex-session"

func NewServer(
	env string,
	config *viper.Viper,
	logger *logrus.Entry,
	assets fs.FS,
	decoder *schema.Decoder,
	hrmConn *grpc.ClientConn,
) (*mux.Router, error) {
	s := &Server{
		env:     env,
		logger:  logger,
		assets:  assets,
		assetFS: hashfs.NewFS(assets),
		decoder: decoder,
		config:  config,
		reg: struct {
			signUpG.SignUpServiceClient
		}{
			SignUpServiceClient: signUpG.NewSignUpServiceClient(hrmConn),
		},
		log: struct {
			logING.LoginServiceClient
		}{
			LoginServiceClient: logING.NewLoginServiceClient(hrmConn),
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
	mw.ChainHTTPMiddleware(r, logger,
		mw.CSRF([]byte(csrfSecret), csrf.Secure(csrfSecure), csrf.Path("/")),
	)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", cacheStaticFiles(http.FileServer(http.FS(s.assetFS)))))

	r.HandleFunc(homePath, s.getHomeHandler).Name("home")

	/* signup */
	r.HandleFunc(signUpPath, s.getSignUpHandler).Methods("GET").Name("signup-form")
	r.HandleFunc(signUpPath, s.postSignUpHandler).Methods("POST").Name("signup")

	/* login */
	r.HandleFunc(loginPath, s.getLoginHandler).Methods("GET").Name("login-form")
	r.HandleFunc(loginPath, s.postLoginHandler).Methods("POST").Name("login")

	r.NotFoundHandler = s.getErrorHandler()
	return r, nil
}
