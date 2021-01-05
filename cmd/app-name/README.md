## Database migration

Database changes are executed using `[golang-migrate](https://github.com/golang-migrate/migrate)`.

## Graceful shutdown

When our server receives an OS signal to shut down we do the following:

- Stop receiving new requests
- Wait for ongoing requests to complete
- Release shared resources like database connections 
