package handler

import (
	"html/template"
	"io/fs"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	user "practice/webex/gunk/v1/user"
)

type Server struct {
	templates *template.Template
	decoder   *schema.Decoder
	sess      *sessions.CookieStore
	user      userSignUp
	env       string
	logger    *logrus.Entry
	assetFS   *hashfs.FS
	assets    fs.FS
}

type userSignUp interface {
	user.UserServiceClient
}
