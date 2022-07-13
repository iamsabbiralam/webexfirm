package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"github.com/yookoala/realpath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"personal/webex/cms/handler"
	"personal/webex/serviceutil/logging"
)

const (
	svcName = "webex"
	version = "1.0.0"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := logging.NewLogger().WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	switch config.GetString("runtime.loglevel") {
	case "trace":
		log.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	default:
		log.Logger.SetLevel(logrus.InfoLevel)
	}
	log.WithField("log level", log.Logger.Level).Info("starting cms service")

	// dialing to hrm api
	u := config.GetString("hrm.url")
	opts := getGRPCOpts(config, false)
	log.Info("dialing hrm api...")
	hrmConn, err := grpc.Dial(u, opts...)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to hrm api")
	}
	defer hrmConn.Close()

	s, err := newServer(log, config, hrmConn)
	if err != nil {
		return err
	}

	s.Use(func(h http.Handler) http.Handler {
		recov := negroni.NewRecovery()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recov.ServeHTTP(w, r, h.ServeHTTP)
		})
	})

	l, err := net.Listen("tcp", ":"+config.GetString("server.port"))
	if err != nil {
		return err
	}

	if err := http.Serve(l, s); err != nil {
		return err
	}
	return nil
}

func newServer(log *logrus.Entry, config *viper.Viper, hrmConn *grpc.ClientConn) (*mux.Router, error) {
	env := config.GetString("runtime.environment")
	log.WithField("environment", env).Info("configuring service")

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		return nil, err
	}
	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))

	srv, err := handler.NewServer(env, config, log, asst, decoder, hrmConn)
	return srv, err
}

func getGRPCOpts(cnf *viper.Viper, withTLS bool) []grpc.DialOption {
	var opts []grpc.DialOption
	// todo(robin): fix issue with not being able to connect to profile with tls
	// rt := cnf.GetString("runtime.environment")
	// if rt == "localdev" || rt == "staging" {
	if withTLS {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	}
	opts = append(opts, grpc.WithBlock())
	return opts
}
