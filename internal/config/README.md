## Configuration managment

All of the application's configuration are present as environment variables.

The application reads these environment variables and sets them to a local `map[string]string` in the `config` package.

The `config` package exposes a function called `func Get() (Config, error)`. This func returns `Config` which contains all the environment variables in golang's data types. For instance the `JWT_TTL` environment variable is stored as `time.Duration`. This function also validates if all the environment variables are present and have a valid value, if not it returns an error.

The following are the supported environment variables:

* PORT: Port at which the server will listen for requests
* MIGRATION_SOURCE_PATH: Path to directory containing the migration files
* SEED_EMAIL_ID: Root administrator email id
* SEED_PHONE_NUMBER: Root administrator phone number
* MONGO_URI: URI to connect to MongoDB
* MONGO_DB_NAME: MongoDB database name
* LOG_LEVEL: Level of logs. Valid value can be found [here](https://github.com/dannypaul/go-skeleton/tree/master/cmd/app-name#logging)
* JWT_SECRET: Secret with which JWT signatures are generated
* JWT_TTL: Time to live(TTL) of JWT
* CHALLENGE_TTL: Time to live(TTL) for a identity challenge
