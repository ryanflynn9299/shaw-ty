package main

import (
	"URLShortener/api/controllers"
	"URLShortener/api/routes"
	"URLShortener/internal/config"
	"URLShortener/internal/i18n"
	"URLShortener/internal/services"
	"URLShortener/internal/storage/db"
	"URLShortener/internal/utils"
	"URLShortener/sql/migrations"
	"URLShortener/sql/seeding"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun/migrate"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate Machine ID - TODO: move this validation logic to load
	if validMId := utils.ValidateMachineId(cfg.MachineID); !validMId {
		log.Fatalf("Invalid machine ID: %d", cfg.MachineID)
	}

	_ = i18n.Load()
	i18n.SetDefaultLocale(cfg.DefaultLocale)

	ctx := context.Background()

	// Initialize the database and run migrations
	dbconn, _ := db.NewDBConn(cfg)
	defer dbconn.Close()
	migrator := migrate.NewMigrator(dbconn, migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize database and migrator: %v", err)
	}

	if group, err := migrator.Migrate(ctx); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	} else if group.IsZero() {
		log.Println("No migrations to run")
	} else {
		log.Printf("Migrated database: %v", group)
	}

	// Run Dataseeding - DEV only
	if cfg.IsDevMode {
		log.Println("Dev Mode: Running Data seeding for testing.")
		err := seeding.SeedUsers(ctx, dbconn)
		if err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
	}

	// Initialize repositories
	shortlinkRepo := db.NewShortLinkRepositoryDB(dbconn)
	userRepo := db.NewUserRepositoryDB(dbconn)

	// Initialize services
	shortlinkService := services.NewLinkService(shortlinkRepo, ctx)
	userService := services.NewUserService(userRepo, ctx)

	// Initialize API endpoint controllers
	authCntlr := controllers.NewAuthController(&userService)
	userCntlr := controllers.NewUserController(&userService)
	linkCntlr := controllers.NewLinkController(&shortlinkService)

	// Setup and start the router
	router := gin.Default()
	routes.InitRouter(router, &userCntlr, &linkCntlr, &authCntlr)
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
