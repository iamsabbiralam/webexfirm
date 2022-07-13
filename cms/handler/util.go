package handler

import (
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

type TemplateData struct {
	Env       string
	CSRFField template.HTML
}

func cacheStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if asset is hashed extend cache to 180 days
		e := `"4FROTHS24N"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=15552000")
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) getErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.doTemplate(w, r, "error.html", http.StatusTemporaryRedirect); err != nil {
			s.logger.WithError(err).Error("unable to load error template")
		}
	})
}

func (s *Server) templateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       s.env,
		CSRFField: csrf.TemplateField(r),
	}
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, name string, status int) error {
	template := s.lookupTemplate(name)
	if template == nil || isPartialTemplate(name) {
		template, status = s.templates.Lookup("error.html"), http.StatusNotFound
	}

	w.WriteHeader(status)
	return template.Execute(w, s.templateData(r))
}

func isPartialTemplate(name string) bool {
	return strings.HasSuffix(name, ".part.html")
}

func (s *Server) lookupTemplate(name string) *template.Template {
	if s.env == "development" {
		if err := s.parseTemplates(); err != nil {
			s.logger.WithError(err).Error("template reload")
			return nil
		}
	}
	return s.templates.Lookup(name)
}

func (s *Server) parseTemplates() error {
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
		"activeStatus": func(status int32) string {
			if status == 1 {
				return "Active"
			}
			return "Inactive"
		},
		"incrementKey": func(status int) int {
			return status + 1
		},
		"formatDate": func(ts *tspb.Timestamp, layout string) string {
			if !ts.IsValid() {
				return ""
			}
			return ts.AsTime().Format(layout)
		},

		"countPaginate": func(a, b int32) int32 {
			if a > 0 {
				c := a / b
				if a%b != 0 {
					c = c + 1
				}
				return c
			}
			return 0
		},
		"noscape": func(str string) template.HTML {
			return template.HTML(str)
		},
		"camelcase": func(str string) string {
			texts := strings.Split(str, "-")
			mainText := []string{}
			for _, v := range texts {
				mainText = append(mainText, strings.Title(v))
			}
			return strings.Join(mainText, "  ")
		},
		"nowtime": func() string {
			return time.Now().Format("02 Jan 2006")
		},
		"permissionChecked": func(res string, act string, allPerm map[string][]string) string {
			val := allPerm[res]
			for _, v := range val {
				if v == act {
					return "checked"
				}
			}
			return ""
		},
		"permission": func(res string) bool {
			return true
		},
		"urls": func(url string, params ...string) string {
			for _, v := range params {
				a := strings.Split(v, "_")
				if len(a) == 2 {
					url = strings.Replace(url, "{"+a[0]+"}", a[1], 1)
				}
			}
			return url
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.assets, "templates/*/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}
