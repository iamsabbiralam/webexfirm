package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	signUpG "personal/webex/gunk/v1/signUp"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type transport struct {
	Transport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rq := req.Clone(req.Context())
	rq.Header.Set("X-Forwarded-Proto", "https")
	return t.Transport.RoundTrip(rq)
}

type Authenticator struct {
	Config    oauth2.Config
	BaseURL   string
	LogoutURL string
}

type Server struct {
	templates *template.Template
	env       string
	logger    *logrus.Entry
	assets    fs.FS
	assetFS   *hashfs.FS
	decoder   *schema.Decoder
	config    *viper.Viper
	reg       signUp
}

type signUp interface {
	signUpG.SignUpServiceClient
}
