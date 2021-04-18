package main

import (
	"fmt"

	"github.com/PSauerborn/geoprowler/pkg/api"
	"github.com/PSauerborn/geoprowler/pkg/utils"
)

var cfg = utils.NewConfigMapWithValues(
	map[string]string{
		"postgres_host":     "192.168.99.100",
		"postgres_port":     "5432",
		"postgres_username": "postgres",
		"postgres_password": "postgres-dev",
		"postgres_database": "geoprowler",
	},
)

func main() {
	// generate new postgres persistence layer and connect
	persistence := api.NewPersistence(cfg.Get("postgres_host"), 5432,
		cfg.Get("postgres_user"), cfg.Get("postgres_password"), cfg.Get("postgres_database"))
	if err := persistence.Connect(); err != nil {
		panic(fmt.Sprintf("unable to generate persistence layer: %+v", err))
	}
	defer persistence.Close()
	// generate new gin-gonic engine/router and start service
	api.NewEngine().Run(fmt.Sprintf(":%d", 10876))
}
