## Configuration managment

All of the application's configuration are present as environment variables.

The application reads these environment variables and sets them to a local `map[string]string`.

The `config` package exposes a function called `func Get() (Config, error)`. This func returns `Config` which contains all the environment variables in golang's data types. For instance the `JWT_TTL` environment variable is stored as `time.Duration` .