package handler

import (
	"text/template"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	user "practice/webex/gunk/v1/user"
)

type Server struct {
	templates *template.Template
	decoder   *schema.Decoder
	sess      *sessions.CookieStore
	user      userSignUp
}

type userSignUp interface {
	user.UserServiceClient
}
