package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"practice/webex/serviceutil/logging"

	usr "practice/webex/gunk/v1/user"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email    string
	Password string
}

type LoginTempData struct {
	CSRFField   template.HTML
	Form        Login
	FormAction  string
	GlobalURLs  map[string]string
	FormErrors  map[string]string
	FormMessage map[string]string
}

func (l Login) Validate(server *Server, r *http.Request, id string) error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, vre("The email field is required"), validation.Length(3, 100), emailValidation(server, r, l.Email)),
		validation.Field(&l.Password, vre("The password field is required"), passwordValidation(server, r, l.Email, l.Password)),
	)
}

func (s *Server) loadLoginForm(w http.ResponseWriter, r *http.Request, data LoginTempData) {
	template := s.lookupTemplate("login.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusSeeOther)
		return
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, "", http.StatusSeeOther)
	}
}

func (s *Server) getLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getLoginHandler")
	data := LoginTempData{
		CSRFField:  csrf.TemplateField(r),
		FormAction: loginURL,
		GlobalURLs: adminViewURLs(),
	}
	s.loadLoginForm(w, r, data)
}

func (s *Server) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postLoginHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "error parsing form"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}

	var form Login
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	if err := form.Validate(s, r, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}

		data := LoginTempData{
			FormErrors: vErrs,
			FormAction: loginURL,
			CSRFField:  csrf.TemplateField(r),
		}

		s.loadLoginForm(w, r, data)
		return
	}

	res, err := s.user.GetUser(ctx, &usr.GetUserRequest{
		User: &usr.User{
			Email:    form.Email,
		},
	})
	if err != nil {
		log.WithError(err).Error("error getting user")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	session, err := s.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Options.HttpOnly = true
	session.Values["authUserID"] = res.User.ID
	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, dashboardPath, http.StatusSeeOther)
}

func emailValidation(s *Server, r *http.Request, email string) validation.Rule {
	return validation.By(func(interface{}) error {
		ctx := r.Context()
		logging.FromContext(r.Context()).WithField("method", "emailValidation")
		_, err := s.user.GetUser(ctx, &usr.GetUserRequest{
			User: &usr.User{
				Email: email,
			},
		})

		if err != nil {
			return fmt.Errorf("invalid email given")
			} else {
			return nil
		}
	})
}

func passwordValidation(s *Server, r *http.Request, email, pass string) validation.Rule {
	return validation.By(func(interface{}) error {
		ctx := r.Context()
		logging.FromContext(r.Context()).WithField("method", "passwordValidation")
		res, _ := s.user.GetUser(ctx, &usr.GetUserRequest{
			User: &usr.User{
				Email:    email,
			},
		})

		if err := bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(pass)); err != nil {
			return fmt.Errorf("invalid password given")
		}

		return nil
	})
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	log := s.logger.WithField("method", "logoutHandler")
	fmt.Println("---session---")
	session, err := s.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Values["authUserID"] = nil
	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, homeURL, http.StatusTemporaryRedirect)
}
