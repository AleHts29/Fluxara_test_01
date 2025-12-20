package main

import (
	"fluxara/internal/adapters/repos/db"
	"fluxara/internal/adapters/rest"
	"fluxara/internal/config"
	serviceDb "fluxara/internal/services/repos/db"
	"log"
	"os"
)

func main() {
	log.Println("Fluxara live!")
	// inicia configs
	config.Load()

	adapterDb, err := db.NewDbAdapter(config.Get())
	if err != nil {
		log.Panic("Error creando el adapter desde main")
		os.Exit(1)
	}
	serviceDb := serviceDb.NewDbService(adapterDb)

	handlers := rest.NewHandlers(serviceDb)
	rest.NewRouter(config.Get(), handlers)
}
