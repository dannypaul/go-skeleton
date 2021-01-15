package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	"github.com/dannypaul/go-skeleton/internal/config"
	"github.com/dannypaul/go-skeleton/internal/driver/platform/mongo"
	"github.com/dannypaul/go-skeleton/internal/iam"
	"github.com/dannypaul/go-skeleton/internal/middleware"
	"github.com/dannypaul/go-skeleton/internal/notification"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading environment variables")
	}

	logLevel, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Error generating log level")
	}
	log.Info().Msg("Log level set to " + logLevel.String())
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	mongoDbClient := mongo.Connect(ctx)

	log.Info().Msg("Starting database migration")

	migrationDriver, err := mongodb.WithInstance(mongoDbClient.Client, &mongodb.Config{DatabaseName: conf.MongoDbName})
	if err != nil {
		log.Fatal().Err(err).Msg("error initialising MongoDB migration driver")
	}

	migration, err := migrate.NewWithDatabaseInstance(conf.MigrationSourcePath, conf.MongoDbName, migrationDriver)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initialising migration")
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Error running migration")
	}

	log.Info().Msg("Successfully completed database migration")

	notificationService := notification.NewService()

	userRepo, _ := iam.NewMongoUserRepo(mongoDbClient)
	challengeRepo, _ := iam.NewMongoChallengeRepo(mongoDbClient)
	iamService := iam.NewService(userRepo, challengeRepo, notificationService)

	_ = iamService.VerifySeedUser(ctx)

	router := chi.NewRouter()

	router.Use(middleware.RequestId, middleware.Auth)

	router.Mount("/identity", iam.Router(iamService))

	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}

	var osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server stopped because of an error")
		}
	}()

	log.Info().Msg("Successfully started server")
	s := <-osSignal
	log.Info().Msgf("Received os signal: %+v", s)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		log.Info().Msg("Server shutdown successful")

		// Release all shared resources
		mongo.Disconnect(mongoDbClient)

		log.Info().Msg("Released all shared resources")

		cancel()

		os.Exit(0)
	}()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Info().Msgf("Server shutdown failed %+v", err)
	}
}
