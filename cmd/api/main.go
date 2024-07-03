package main

import (
	"log"

	_ "cleanarch/boiler/internal/docs"
	"cleanarch/boiler/internal/server"
)

func main() {

	app := server.NewApp()

	if err := app.Run("3000"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
