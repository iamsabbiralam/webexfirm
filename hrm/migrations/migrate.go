package main

import (
	"log"
	"practice/webex/hrm/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
