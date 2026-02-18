package main

import (
	"fluxara/internal/adapters/repos/db"
	mdpadapter "fluxara/internal/adapters/repos/mp"
	"fluxara/internal/adapters/rest"
	"fluxara/internal/config"
	serviceDb "fluxara/internal/services/repos/db"
	serviceDbGergal "fluxara/internal/services/repos/dbGergal"
	serviceMp "fluxara/internal/services/repos/mp"
	"fmt"
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

	adapterMp, err := mdpadapter.NewMPAdapter(config.Get())
	serviceMp := serviceMp.NewDbService(adapterDb, adapterMp)

	adapterDbGergal, err := db.NewDbAdapter(config.Get())
	if err != nil {
		log.Panic("Error creando el adapter Gergal desde main")
		os.Exit(2)
	}

	serviceDbGergal := serviceDbGergal.NewDbServiceGergal(adapterDbGergal)

	handlers := rest.NewHandlers(serviceDb, serviceDbGergal, serviceMp)
	rest.NewRouter(config.Get(), handlers)

	fmt.Println("Esto es config", config.Get())
}
