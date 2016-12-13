package main

import (
	"log"

	"github.com/byuoitav/authmiddleware/bearertoken"
)

func main() {
	log.Println("Start")

	token, err := bearertoken.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", token)
}
