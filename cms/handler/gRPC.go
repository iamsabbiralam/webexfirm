package handler

import (
	"html/template"
	"io/fs"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	ccG "practice/webex/gunk/v1/circularCategory"
	user "practice/webex/gunk/v1/user"
)

type Server struct {
	templates *template.Template
	decoder   *schema.Decoder
	sess      *sessions.CookieStore
	env       string
	logger    *logrus.Entry
	assetFS   *hashfs.FS
	assets    fs.FS
	user      userSignUp
	cc        circularCategory
}

type userSignUp interface {
	user.UserServiceClient
}

type circularCategory interface {
	ccG.CircularCategoryServiceClient
}
