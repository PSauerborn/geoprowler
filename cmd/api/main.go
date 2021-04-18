package main

import (
	"fmt"

	"github.com/PSauerborn/geoprowler/pkg/api"
)

func main() {
	// generate new postgres persistence layer and connect
	persistence := api.NewPersistence("192.168.99.100", 5432,
		"postgres", "postgres-dev", "geoprowler")
	if err := persistence.Connect(); err != nil {
		panic(fmt.Sprintf("unable to generate persistence layer: %+v", err))
	}
	defer persistence.Close()
	// generate new gin-gonic engine/router and start service
	api.NewEngine().Run(fmt.Sprintf(":%d", 10999))
}
