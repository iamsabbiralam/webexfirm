package handler

import (
	"errors"
	"html/template"
	"net/http"
	"personal/webex/serviceutil/logging"
	"unicode"

	signG "personal/webex/gunk/v1/signUp"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SignUpForm struct {
	CSRFField            template.HTML
	FormAction           string
	FirstName            string
	LastName             string
	Username             string
	Email                string
	Password             string
	PasswordConfirmation string
	FormErrors           map[string]error
}

func (s *Server) getSignUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "getSignUpHandler")
	form := SignUpForm{
		CSRFField:  csrf.TemplateField(r),
		FormAction: signUpPath,
	}
	template := s.templates.Lookup("sign-up.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, form); err != nil {
		errMsg := "error with template execution"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (s *Server) postSignUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postSignUpHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}

	var form SignUpForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	formError := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.FirstName, validation.Required, validation.Length(3, 64)),
		validation.Field(&form.LastName, validation.Required, validation.Length(3, 64)),
		validation.Field(&form.Email, validation.Required, is.Email),
		validation.Field(&form.Username, validation.Required, validation.Length(3, 64)),
		validation.Field(&form.Password, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.PasswordConfirmation, validation.Required, validation.Length(8, 64)),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Email"] != nil {
				formError["Email"] = err["Email"]
			}

			if err["Password"] != nil {
				formError["Password"] = err["Password"]
			}

			if err["Phone"] != nil {
				formError["Phone"] = err["Phone"]
			}

			if err["PasswordConfirmation"] != nil {
				formError["PasswordConfirmation"] = err["PasswordConfirmation"]
			}

			if err["FirstName"] != nil {
				formError["FirstName"] = err["FirstName"]
			}

			if err["LastName"] != nil {
				formError["LastName"] = err["LastName"]
			}

			if err["Username"] != nil {
				formError["Username"] = err["Username"]
			}
		}

		logging.WithError(err, log).Error("invalid request")
	}

	if formError["Password"] == nil && formError["PasswordConfirmation"] == nil && !isPasswordValid(form.Password) {
		formError["Password"] = errors.New("Password must be at least 8 characters long, contain at least one upper case letter, one lower case letter, one number and one special character")
	}

	if formError["Password"] == nil && formError["PasswordConfirmation"] == nil && form.Password != form.PasswordConfirmation {
		formError["PasswordConfirmation"] = errors.New("Password and Password Confirmation must match")
	}

	form.FormErrors = formError
	if len(formError) > 0 {
		form.CSRFField = csrf.TemplateField(r)
		template := s.templates.Lookup("sign-up.html")
		if template == nil {
			http.Redirect(w, r, "unable to load template", http.StatusSeeOther)
			return
		}

		if err := template.Execute(w, form); err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
		}
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.WithError(err, log).Error("error with encrypted password")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	_, err = s.reg.Registration(ctx, &signG.RegisterRequest{
		SignUP: &signG.SignUP{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Username:  form.Username,
			Email:     form.Email,
			Password:  string(pass),
			Status:    signG.Status_ACTIVE,
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("error with registration")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func isPasswordValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
