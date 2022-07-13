package main

import (
	"log"

	"personal/webex/hrm/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
