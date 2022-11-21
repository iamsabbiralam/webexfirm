package mw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"practice/webex/serviceutil/logging"
)

func ChainHTTPMiddleware(usr *mux.Router, log logrus.FieldLogger, mw ...func(http.Handler) http.Handler) {
	for _, f := range []func(http.Handler) http.Handler{
		Logger(log),
		Gzip,
		ContentType("text/html"),
	} {
		usr.Use(f)
	}
	for _, f := range mw {
		usr.Use(f)
	}
}

func Logger(log logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			lrw := negroni.NewResponseWriter(w)
			ctx := logging.WithLogger(req.Context(), log)
			h.ServeHTTP(lrw, req.WithContext(ctx))

			// TODO: after a while log just errors, for now log everything
			status := lrw.Status()
			ss := lrw.Header()
			hostname := strings.ToLower(req.Host)
			method := req.Method
			path := req.URL.String()
			location := ss.Get("Location")
			ip := GetIP(req)

			log := log.WithFields(logrus.Fields{
				"status":     status,
				"statusText": http.StatusText(status),
				"hostname":   hostname,
				"method":     method,
				"path":       path,
				"clientIP":   ip,
			})
			if location != "" {
				log = log.WithField("redirectLocation", location)
			}

			log.Infof("request handled")
		})
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-forwarded-for")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func CSRF(secret []byte, opts ...csrf.Option) func(h http.Handler) http.Handler {
	opts = append([]csrf.Option{
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.FromContext(r.Context())
			log.WithFields(logrus.Fields{
				"csrf_error": csrf.FailureReason(r).Error(),
				"token":      csrf.Token(r),
				"template":   csrf.TemplateField(r),
			}).Error("csrf error")
			fmt.Fprintln(w, csrf.FailureReason(r))
		})),
	}, opts...)
	return csrf.Protect(secret, opts...)
}

type contentWriter struct {
	http.ResponseWriter
	def string
	log *logrus.Entry
}

func (c contentWriter) Write(b []byte) (int, error) {
	if c.Header().Get("Content-Type") == "" {
		ct := http.DetectContentType(b)
		if ct == "" {
			ct = c.def
		}
		c.log.WithField("content-type", ct).Trace("set")
		c.Header().Set("Content-Type", ct)
	}
	c.log.WithField("content-type", c.Header().Get("Content-Type")).Trace("write")
	return c.ResponseWriter.Write(b)
}

func ContentType(typeDefault string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(contentWriter{
				ResponseWriter: w,
				def:            typeDefault, log: logging.FromContext(r.Context()),
			}, r)
		})
	}
}

func Recovery() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		recov := negroni.NewRecovery()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recov.ServeHTTP(w, r, h.ServeHTTP)
		})
	}
}
