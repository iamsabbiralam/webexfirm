package handler

import (
	"html/template"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
)

type TemplateData struct {
	Env       string
	CSRFField template.HTML
}

type DynamicQueryString struct {
	SearchTerm  string
	PageNumber  int32
	CurrentPage int32
	Offset      int32
	OthersValue map[string]string
}

func GetQueryStringData(r *http.Request, keys []string, isNotDefault bool) *DynamicQueryString {
	var data DynamicQueryString
	queryParams := r.URL.Query()
	var err error
	if !isNotDefault {
		data.SearchTerm, err = url.PathUnescape(queryParams.Get("search-term"))
		if err != nil {
			data.SearchTerm = ""
		}
		page, err := url.PathUnescape(queryParams.Get("page"))
		if err != nil {
			page = "1"
		}
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			pageNumber = 1
		}
		data.PageNumber = int32(pageNumber)
		var offset int32 = 0
		currentPage := pageNumber
		if currentPage <= 0 {
			currentPage = 1
		} else {
			offset = limitPerPage*int32(currentPage) - limitPerPage
		}
		data.CurrentPage = int32(currentPage)
		data.Offset = offset
	}
	if len(keys) > 0 {
		for _, v := range keys {
			data.OthersValue[v], err = url.PathUnescape(queryParams.Get(v))
			if err != nil {
				data.OthersValue[v] = ""
			}
		}
	}
	return &data
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

func (s *Server) parseTemplates() error {
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.assets, "templates/*/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

// Round float to 2 decimal places
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}