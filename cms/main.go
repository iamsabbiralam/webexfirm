package main

import (
	"fmt"
	"practice/webex/cms/handler"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	user "practice/webex/gunk/v1/user"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	
	store := sessions.NewCookieStore([]byte(config.GetString("session.key")))

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.GetString("categoryService.host"), config.GetString("categoryService.port")),
	grpc.WithInsecure(),
)
	if err != nil {
		log.Fatal(err)
	}
	tc := user.NewUserServiceClient(conn)
	r := handler.New(decoder, store, tc)

	host, port := config.GetString("server.host"), config.GetString("server.port")

	log.Printf("Server starting on http://%s:%s", host, port)
	if err:= http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}

}