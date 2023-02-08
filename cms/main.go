package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"practice/webex/cms/handler"
	"practice/webex/serviceutil/logging"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/yookoala/realpath"
	"google.golang.org/grpc"
)

const (
	svcName = "webex-website"
	version = "1.0.0"
)

func main() {
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
		log.Printf("error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	store := sessions.NewCookieStore([]byte(config.GetString("session.key")))
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		log.Fatal(err)
	}

	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.GetString("userService.host"), config.GetString("userService.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	r, _ := handler.Handler(decoder, config, log, store, conn, asst)
	host, port := config.GetString("server.host"), config.GetString("server.port")
	log.Printf("Server starting on http://%s:%s", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
