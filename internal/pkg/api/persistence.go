package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/geoprowler/internal/pkg/utils"
)

var (
	// define custom errors
	ErrEntityNotFound    = errors.New("cannot find entity with provided entity ID")
	ErrInvalidEntityMeta = errors.New("entity has invalid metadata")
)

type Persistence struct {
	*utils.PostgresPersistence
}

// function used to generate postgres connection
// string from input parameters
func generateConnectionString(host string, port int, username,
	password, database string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password,
		host, port, database)
}

// function used to generate new persistence
// instance to retrieve and manage data
func NewPersistence(host string, port int, username,
	password, database string) *Persistence {
	// generate new instance of base postgres
	// persistence
	base := utils.PostgresPersistence{
		DatabaseURL: generateConnectionString(host, port, username, password, database),
	}
	return &Persistence{
		PostgresPersistence: &base,
	}
}

type Entity struct {
	EntityId    uuid.UUID              `json:"entity_id"`
	LastUpdated *time.Time             `json:"last_updated"`
	Location    *GeoLocation           `json:"location"`
	Meta        map[string]interface{} `json:"meta"`
}

// function to retrieve single entity from the postgres database
func (db *Persistence) GetEntity(entityId uuid.UUID) (Entity, error) {
	log.Debug(fmt.Sprintf("fetching entity %s from database...", entityId))
	var (
		e        Entity
		location map[string]float64
	)

	query := `SELECT entities.entity_id, entities.meta, locations.last_updated,
	locations.location FROM entities INNER JOIN locations ON entities.entity_id
	= locations.entity_id WHERE entities.entity_id = $1`
	row := db.Session.QueryRow(context.Background(), query, entityId)
	if err := row.Scan(&e.EntityId, &e.Meta, &e.LastUpdated, &location); err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return e, ErrEntityNotFound
		default:
			return e, err
		}
	}

	// convert location map into JSON string and convert to struct
	jsonBody, _ := json.Marshal(location)
	if err := json.Unmarshal(jsonBody, &e.Location); err != nil {
		log.Error(fmt.Errorf("retrieved invalid geolocation %s: %+v", jsonBody, err))
		return e, err
	}
	return e, nil
}

// function to retrieve all entities from the postgres database
func (db *Persistence) GetEntities() ([]Entity, error) {
	log.Debug("fetching all entities from database...")
	entities := []Entity{}

	query := `SELECT entities.entity_id, entities.meta, locations.last_updated,
	locations.location FROM entities INNER JOIN locations ON entities.entity_id
	= locations.entity_id`
	results, err := db.Session.Query(context.Background(), query)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return entities, nil
		default:
			return entities, err
		}
	}

	for results.Next() {
		var (
			e        Entity
			location map[string]float64
		)
		// scan values into location variables
		if err := results.Scan(&e.EntityId, &e.Meta, &e.LastUpdated, &location); err != nil {
			log.Error(fmt.Errorf("unable to scan data into local variables: %+v", err))
			continue
		}
		// convert location map into JSON string and convert to struct
		jsonBody, _ := json.Marshal(location)
		if err := json.Unmarshal(jsonBody, &e.Location); err != nil {
			log.Error(fmt.Errorf("retrieved invalid geolocation %s: %+v", jsonBody, err))
			continue
		}
		// add entity to list of entities
		entities = append(entities, e)
	}
	return entities, nil
}

// function to add new entity to the postgres database
func (db *Persistence) RegisterEntity(meta map[string]interface{}) error {
	log.Debug("inserting new entity into the database...")
	jsonBody, err := json.Marshal(meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert metadata to JSON: %+v", err))
		return ErrInvalidEntityMeta
	}

	var query string
	entityId := uuid.New()
	// insert entry into entities table
	query = `INSERT INTO entities(entity_id, meta) VALUES($1, $2)`
	_, err = db.Session.Exec(context.Background(), query, entityId, jsonBody)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		return err
	}
	// insert entry into locations table
	query = `INSERT INTO locations(entity_id) VALUES($1)`
	_, err = db.Session.Exec(context.Background(), query, entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		return err
	}
	return nil
}

// function used to delete a given entity
func (db *Persistence) DeleteEntity(entityId uuid.UUID) error {
	log.Debug(fmt.Sprintf("deleting entity %s from datbase...", entityId))
	var (
		query string
		err   error
	)
	// delete entry from entities table
	query = `DELETE FROM entities WHERE entity_id = $1`
	_, err = db.Session.Exec(context.Background(), query, entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		return err
	}
	// delete entries from locations table
	query = `DELETE FROM locations WHERE entity_id = $1`
	_, err = db.Session.Exec(context.Background(), query, entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		return err
	}
	return nil
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// set function on GeoLocation instance to check for validity
func (l GeoLocation) IsValid() bool {
	// define condition for latitude
	validLatitude := l.Latitude >= -90 && l.Latitude <= 90
	if !validLatitude {
		log.Warn(fmt.Sprintf("latitude %f invalid", l.Latitude))
	}
	// define condition for longitude
	validLongitude := l.Longitude >= -180 && l.Longitude <= 180
	if !validLongitude {
		log.Warn(fmt.Sprintf("longitude %f invalid", l.Longitude))
	}
	return validLatitude && validLongitude
}

// function used to set location for a given entity
func (db *Persistence) SetEntityLocation(entityId uuid.UUID,
	location GeoLocation) error {
	log.Debug(fmt.Sprintf("setting location for entity %s...", entityId))

	jsonBody, _ := json.Marshal(location)
	query := `UPDATE locations SET location = $1, last_updated = $2
	WHERE entity_id = $3`

	_, err := db.Session.Exec(context.Background(), query, jsonBody, time.Now().UTC(), entityId)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute database query: %+v", err))
		return err
	}
	return nil
}
