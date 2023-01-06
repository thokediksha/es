package main

import (
	Routers "es/routers"
	"log"
)

func main() {
	router := Routers.SetupRouter()

	if err := router.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}
