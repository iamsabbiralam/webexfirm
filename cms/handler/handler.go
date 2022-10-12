package handler

import (

	// "log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	user "practice/webex/gunk/v1/user"
)

const sessionName = "cms-session"

type Handler struct {
	templates *template.Template
	decoder *schema.Decoder
	sess *sessions.CookieStore
	tc	user.UserServiceClient
}

func New(decoder *schema.Decoder, sess *sessions.CookieStore, tc user.UserServiceClient) *mux.Router {
	h:= &Handler{
		decoder: decoder,
		sess: sess,
		tc: tc,
	}

	h.parseTemplate()

	r:= mux.NewRouter()
	r.HandleFunc("/", h.home)
	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))
	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, "invalid URL", http.StatusInternalServerError)
			return
		}
	})

	return r
}

func (h *Handler) parseTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"cms/assets/templates/user/create-user.html",
		"cms/assets/templates/home.html",
	))
}
