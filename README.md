# geoprowler

`GeoProwler` is a containerized, stateless REST interface for location
tracking and management, backed and persisted in a `Postgres` instance.
The current feature set allows users to generate entities (with optional
JSON metadata for flexibility and integration). The location of entities
can then be updated via a `PUT` request to the `/location/{entityId}`
endpoint.

See `docs/swagger.yaml` for full documentation on API routes
and REST interface, and `docs/db.sql` for a `Postgres` dump of the required
tables

