package handler

import (
	"html/template"
	"log"
	"net/http"
	"practice/webex/serviceutil/logging"
	"time"

	usr "practice/webex/gunk/v1/user"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type SignUp struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	UserName    string
	Gender      string
	PhoneNumber string
	DOB         time.Time
}

type SignUpTempData struct {
	CSRFField   template.HTML
	Form        SignUp
	FormAction  string
	FormErrors  map[string]string
	FormMessage map[string]string
}

func (sign SignUp) Validate(server *Server, r *http.Request, id string) error {
	vre := validation.Required.Error
	return validation.ValidateStruct(&sign,
		validation.Field(&sign.FirstName, vre("The First Name is required")),
		validation.Field(&sign.LastName, vre("The Last Name is required")),
		validation.Field(&sign.Email, vre("The Email is required")),
		validation.Field(&sign.Password, vre("The Password is required")),
	)
}

func (s *Server) loadSignUpForm(w http.ResponseWriter, r *http.Request, data SignUpTempData) {
	template := s.lookupTemplate("create-user.html")
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

func (s *Server) signUpMethod(w http.ResponseWriter, r *http.Request) {
	data := SignUpTempData{
		CSRFField:   csrf.TemplateField(r),
		FormAction:  registrationURL,
		FormErrors:  map[string]string{},
		FormMessage: map[string]string{},
	}

	s.loadSignUpForm(w, r, data)
}

func (s *Server) postSignUpMethod(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "postSignUpMethod")
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form SignUp
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, registrationURL, http.StatusSeeOther)
	}

	if err := form.Validate(s, r, ""); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			data := SignUpTempData{
				CSRFField:   csrf.TemplateField(r),
				Form:        form,
				FormAction:  registrationURL,
				FormErrors:  vErrs,
				FormMessage: map[string]string{},
			}
			s.loadSignUpForm(w, r, data)
			return
		}
		http.Redirect(w, r, registrationURL, http.StatusSeeOther)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.WithError(err, log).Error("error with encrypted password")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	_, err = s.user.CreateUser(ctx, &usr.CreateUserRequest{
		User: &usr.User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  string(pass),
			Status:    usr.Status_Inactive,
		},
	})
	if err != nil {
		log.Infof("error while creating user: %+v", err)
		http.Redirect(w, r, registrationURL, http.StatusSeeOther)
	}

	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}
