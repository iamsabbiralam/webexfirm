package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"personal/webex/hrm/storage/postgres"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
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

	store, err := newDBFromConfig(config)
	if err != nil {
		log.Print("unable to configure storage", err)
	}

	if err := setupGRPCService(store, config); err != nil {
		log.Printf("unable to setup grpc service:%+v", err)
	}
}

// NewDBFromConfig build database connection from config file.
func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, db.RunMigration(cf("migrationDir"))
}

func setupGRPCService(store *postgres.Storage, config *viper.Viper) error {
	if err := store.RunMigration(config.GetString("database.migrationDir")); err != nil {
		log.Printf("unable to run migrations")
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GetString("server.port")))
	if err != nil {
		log.Printf("Failed to listen on port 50051: %v", err)
	}
	grpcServer := grpc.NewServer()

	log.Printf("Server hrm management listening at : %+v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Print("Failed to serve GRPC over port : 50051")
		return err
	}
	return nil
}
