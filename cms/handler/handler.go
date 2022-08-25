package handler

import (
	"errors"
	"io/fs"
	"net/http"

	signUpG "personal/webex/gunk/v1/signUp"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

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
	}

	if err := s.parseTemplates(); err != nil {
		return nil, err
	}

	csrfSecret := config.GetString("csrf.secret")
	if csrfSecret == "" {
		return nil, errors.New("CSRF secret must not be empty")
	}

	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", cacheStaticFiles(http.FileServer(http.FS(s.assetFS)))))

	r.HandleFunc(homePath, s.getHomeHandler).Name("home")

	/* signup */
	r.HandleFunc(signUpPath, s.getSignUpHandler).Methods("GET").Name("signup-form")
	r.HandleFunc(signUpPath, s.postSignUpHandler).Methods("POST").Name("signup")

	r.NotFoundHandler = s.getErrorHandler()
	return r, nil
}
