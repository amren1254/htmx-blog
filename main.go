package main

import (
	"log"

	"github.com/amren1254/htmx-blog/router"
	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(".envrc"); err != nil {
		log.Println("error loading envrc file", err)
		return
	}
}

func main() {
	router.InitRoute()
}
