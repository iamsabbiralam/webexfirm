package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"personal/webex/serviceutil/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"

	loginG "personal/webex/gunk/v1/login"
)

type LoginForm struct {
	Email      string
	Password   string
	FormErrors map[string]string
	FormAction string
	CSRFField  template.HTML
}

func (l LoginForm) Validate(server *Server, r *http.Request, id string) error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, vre("The email field is required"), validation.Length(3, 100), emailValidation(server, r, l.Email)),
		validation.Field(&l.Password, vre("The password field is required"), passwordValidation(server, r, l.Email, l.Password)),
	)
}

func (s *Server) loadLoginForm(w http.ResponseWriter, r *http.Request, data LoginForm) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "getLoginHandler")
	template := s.templates.Lookup("login.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, data); err != nil {
		errMsg := "error with template execution"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (s *Server) getLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.FromContext(ctx).WithField("method", "getLoginHandler")
	data := LoginForm{
		CSRFField:  csrf.TemplateField(r),
		FormAction: loginPath,
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

	var form LoginForm
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
		
		data := LoginForm{
			FormErrors: vErrs,
			FormAction: loginPath,
			CSRFField:  csrf.TemplateField(r),
		}
		s.loadLoginForm(w, r, data)
		return
	}

	_, err := s.log.Login(ctx, &loginG.LoginRequest{
		Email: form.Email,
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, dashboardPath, http.StatusTemporaryRedirect)
}

func emailValidation(s *Server, r *http.Request, email string) validation.Rule {
	return validation.By(func(interface{}) error {
		ctx := r.Context()
		logging.FromContext(r.Context()).WithField("method", "emailValidation")
		res, _ := s.log.Login(ctx, &loginG.LoginRequest{
			Email: email,
		})

		if res.Login.Email == email {
			return nil
		} else {
			fmt.Println("email not found")
		}

		return nil
	})
}

func passwordValidation(s *Server, r *http.Request, email, pass string) validation.Rule {
	return validation.By(func(interface{}) error {
		ctx := r.Context()
		logging.FromContext(r.Context()).WithField("method", "emailValidation")
		res, _ := s.log.Login(ctx, &loginG.LoginRequest{
			Email: email,
		})

		if err := bcrypt.CompareHashAndPassword([]byte(res.Login.Password), []byte(pass)); err != nil {
			return fmt.Errorf("invalid password given")
		}

		return nil
	})
}
