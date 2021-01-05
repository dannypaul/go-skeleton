## Database migration

Database changes are executed using [golang-migrate](https://github.com/golang-migrate/migrate).

Each migration consists of two files which have the following naming convention:
- `<index>_<name>.up.<json|sql>`
- `<index>_<name>.down.<json|sql>`

The `up` file should describe the transformation you’d like to make to your database schema, and the `down` file should revert the transformations done by the `up` file. In other words, the database schema should be unchanged if you do an `up` followed by a `down`. For example, if you create a table in the `up` file, you should drop it in the `down` file. It is wise to perform the transformations in precisely the reverse order they were made in the `up` file.

The file extension of the migration file is `json` for MongoDB migrations
Each MongoDB migration file should contain a valid [MongoDB database command](https://docs.mongodb.com/manual/reference/command)

The migration files can be found at `internal/migration`

The path where the migration files are located can be configured using the environment variable `MIGRATION_SOURCE_PATH`.\
Multi stage Docker images typically contain the compiled go binary and other supporting files to run that binary. Migration files is an example of such supporting files. The migration files are mounted in a desired path, and that path is set as the above environment variable. 


## Graceful shutdown

When our server receives an OS signal to shut down we do the following:

- Stop receiving new requests
- Wait for ongoing requests to complete
- Release shared resources like database connections 
