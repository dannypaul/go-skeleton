package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dannypaul/go-skeleton/config"
	"github.com/dannypaul/go-skeleton/internal/driver/platform/mongo"
	"github.com/dannypaul/go-skeleton/internal/iam"
	"github.com/dannypaul/go-skeleton/internal/middleware"
	"github.com/dannypaul/go-skeleton/internal/notification"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		log.Fatal(fmt.Errorf("error reading environment variables %w", err))
	}

	ctx := context.Background()
	mongoDbClient := mongo.Connect(ctx)

	migrationDriver, err := mongodb.WithInstance(mongoDbClient.Client, &mongodb.Config{DatabaseName: conf.MongoDbName})
	if err != nil {
		log.Fatal(fmt.Errorf("error initialising MongoDB migration driver %w", err))
	}

	migration, err := migrate.NewWithDatabaseInstance(conf.MigrationPath, conf.MongoDbName, migrationDriver)
	if err != nil {
		log.Fatal(fmt.Errorf("error initialising migration %w", err))
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(fmt.Errorf("error running migration %w", err))
	}

	notificationService := notification.NewService()

	userRepo, _ := iam.NewMongoUserRepo(mongoDbClient)
	challengeRepo, _ := iam.NewMongoChallengeRepo(mongoDbClient)
	iamService := iam.NewService(userRepo, challengeRepo, notificationService)

	_ = iamService.VerifySeedUser(ctx)

	router := chi.NewRouter()

	router.Use(middleware.RequestIdMiddleware, middleware.AuthMiddleware)

	router.Mount("/identity", iam.Router(iamService))

	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped with error: %+v", err)
		}
	}()

	log.Println("Successfully started server")
	s := <-gracefulStop
	log.Printf("Received signal: %+v", s)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		log.Print("Server shutdown successful")

		// Release all shared resources
		mongo.Disconnect(mongoDbClient)

		log.Print("Released all shared resources")

		cancel()

		os.Exit(0)
	}()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Server shutdown failed %+v", err)
	}
}
