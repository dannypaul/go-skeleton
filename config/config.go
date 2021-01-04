package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type env map[string]string

var e = make(env)

func (e *env) lookup(key string) string {
	if (*e)[key] == "" {
		(*e)[key] = os.Getenv(key)
	}
	return (*e)[key]
}

func (e env) emptyKeys() []string {
	var keys []string
	for key, value := range e {
		if value == "" {
			keys = append(keys, key)
		}
	}
	return keys
}

type Config struct {
	Port            string
	MigrationPath   string
	SeedEmailId     string
	SeedPhoneNumber string
	MongoURI        string
	MongoDbName     string
	JwtSecret       string
	JwtTTL          time.Duration
	ChallengeTTL    time.Duration
}

// Get creates the config by reading the environment variables
func Get() (Config, error) {
	var conf Config

	conf.Port = e.lookup("PORT")

	conf.MigrationPath = e.lookup("MIGRATION_PATH")

	conf.SeedEmailId = e.lookup("SEED_EMAIL_ID")
	conf.SeedPhoneNumber = e.lookup("SEED_PHONE_NUMBER")

	conf.MongoURI = e.lookup("MONGO_URI")
	conf.MongoDbName = e.lookup("MONGO_DB_NAME")

	conf.JwtSecret = e.lookup("JWT_SECRET")
	jwtTTL, err := time.ParseDuration(e.lookup("JWT_TTL"))
	if err != nil {
		return Config{}, err
	}
	conf.JwtTTL = jwtTTL

	challengeTTL, err := time.ParseDuration(e.lookup("CHALLENGE_TTL"))
	if err != nil {
		return Config{}, err
	}
	conf.ChallengeTTL = challengeTTL

	emptyKeys := e.emptyKeys()
	if len(emptyKeys) > 0 {
		err := errors.New("the following environment variables are not set:" + strings.Join(emptyKeys[:], ","))
		return Config{}, err
	}

	return conf, nil
}
