package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.assistagro.com/back/back.auth.go/pkg/aclient"
)

func main() {

	// Create new client
	client, err := aclient.NewClient(os.Getenv("AUTH_URL"), 0)
	if err != nil {
		log.Fatal("Error while creating client: ", err)
	}

	// Get new user session
	session, err := client.NewSession(os.Getenv("AUTH_TOKEN"))
	if err != nil {
		log.Fatal("Error while creating session: ", err)
	}

	user, err := session.GetUser()
	if err != nil {
		log.Fatal("Error while getting user: ", err)
	}
	fmt.Printf("%+v\n", user)
}
