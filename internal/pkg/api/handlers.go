package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var persistence *Persistence

// function used to set global persistence layer
func SetPersistence(host string, port int, username,
	password, database string) *Persistence {
	// generate new persistence instance
	persistence = NewPersistence(host, port, username, password,
		database)
	return persistence
}

// API handler used to serve health check request
func HealthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check handler")
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Service running"})
}

// API handler used to retreive location for a given entity
func GetEntityHandler(ctx *gin.Context) {
	log.Info("received request for location")
	// retrieve entity ID from URL params and parse into UUID
	entityId, err := uuid.Parse(ctx.Param("entityId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse entity ID: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid entity ID format"})
		return
	}

	// get entity (including location) from database
	entity, err := persistence.GetEntity(entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve entity: %+v", err))
		switch err {
		case ErrEntityNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find specified entity"})
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				InternalServerErrorResponse)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK, "entity": entity})
}

// API handler used to retrieve list of current entities
func GetEntitiesHandler(ctx *gin.Context) {
	log.Info("received request to retrieve entities")
	// get entities from database and return
	entities, err := persistence.GetEntities()
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve entities from database: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			InternalServerErrorResponse)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"entities": entities})
}

// API handler used to register new entity
func RegisterEntityHandler(ctx *gin.Context) {
	log.Info("received request to register new entity")

	var e struct {
		Meta map[string]interface{}
	}
	if err := ctx.ShouldBind(&e); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid request body"})
	}

	if err := persistence.RegisterEntity(e.Meta); err != nil {
		log.Error(fmt.Errorf("unable to register new entity: %+v", err))
		switch err {
		case ErrInvalidEntityMeta:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
				"message": "Invalid entity metadata"})
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				InternalServerErrorResponse)
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
		"message": "Successfully created entity"})
}

// API handler used to register new location for a given entity
func RegisterLocationHandler(ctx *gin.Context) {
	log.Info("received request to register new location")
	// retrieve entity ID from URL params and parse into UUID
	entityId, err := uuid.Parse(ctx.Param("entityId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse entity ID: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid entity ID format"})
		return
	}

	var l GeoLocation
	if err := ctx.ShouldBind(&l); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid location"})
		return
	}
	// return 400 if GeoLocation is not valid i.e. longitude and latitude do
	// not fall in allowed ranges
	if !l.IsValid() {
		log.Error(fmt.Sprintf("cannot register location: received invalid coordinates %+v", l))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid location"})
		return
	}

	// get entity from database. if entity cannot be found, return 404
	_, err = persistence.GetEntity(entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve entity: %+v", err))
		switch err {
		case ErrEntityNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find specified entity"})
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				InternalServerErrorResponse)
		}
		return
	}

	// set entity location in database
	if err := persistence.SetEntityLocation(entityId, l); err != nil {
		log.Error(fmt.Errorf("unable to set location for entity: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			InternalServerErrorResponse)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully updated location"})
}

// API handler used to delete a given entity
func DeleteEntityHandler(ctx *gin.Context) {
	log.Info("received request to delete entity")
	// retrieve entity ID from URL params and parse into UUID
	entityId, err := uuid.Parse(ctx.Param("entityId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse entity ID: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"message": "Invalid entity ID format"})
		return
	}

	// get entity from database. if entity cannot be found, return 404
	_, err = persistence.GetEntity(entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve entity: %+v", err))
		switch err {
		case ErrEntityNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"http_code": http.StatusNotFound,
				"message": "Cannot find specified entity"})
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				InternalServerErrorResponse)
		}
		return
	}

	// delete entity from database
	if err := persistence.DeleteEntity(entityId); err != nil {
		log.Error(fmt.Errorf("unable to delete entity: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			InternalServerErrorResponse)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully deleted entity"})
}
