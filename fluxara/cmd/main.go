package main

import (
	"fluxara/internal/adapters/rest"
	"fluxara/internal/config"
	"log"
)

func main() {
	log.Println("Fluxara live!")
	// inicia configs
	config.Load()

	// adapterDb, err := db.NewDbAdapter(config.Get())
	// if err != nil {
	// 	log.Panic("Error creando el adapter desde main")
	// 	os.Exit(1)
	// }
	// serviceDb := serviceDb.NewDbService(adapterDb)

	handlers := rest.NewHandlers(nil)
	rest.NewRouter(config.Get(), handlers)
}
