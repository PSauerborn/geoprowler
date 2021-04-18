package api

import (
	"github.com/gin-gonic/gin"

	"github.com/PSauerborn/geoprowler/internal/pkg/api"
)

// function to generate new API instance. note
// that the persistence instance must be passed
// down to be set as a global variable
func NewEngine() *gin.Engine {
	// generate new instance of gin engine
	// and add routes for REST API
	r := gin.Default()

	r.GET("/health_check", api.HealthCheckHandler)
	// define routes used to retrieve data from API
	r.GET("/entities/:entityId", api.GetEntityHandler)
	r.GET("/entities/all", api.GetEntitiesHandler)

	// define routes to manage entities and update locations
	r.POST("/entities/new", api.RegisterEntityHandler)
	r.PUT("/location/:entityId", api.RegisterLocationHandler)
	r.DELETE("/entities/:entityId", api.DeleteEntityHandler)

	return r
}

// function used to generate new persistence instance
func NewPersistence(host string, port int, username,
	password, database string) *api.Persistence {
	return api.SetPersistence(host, port, username, password,
		database)
}
