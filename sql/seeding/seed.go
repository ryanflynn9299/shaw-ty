package seeding

import (
	"URLShortener/internal/storage/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"log"
)

func SeedUsers(ctx context.Context, db *bun.DB) error {
	usersToSeed := []models.User{
		//{Username: "testuser1", Email: "test1@example.com", HashedPassword: hashPassword("password123")}, // Use a proper hashing func
		//{Username: "testuser2", Email: "test2@example.com", HashedPassword: hashPassword("password456")},
	}

	for _, u := range usersToSeed {
		var existingUser models.User
		err := db.NewSelect().Model(&existingUser).Where("username = ?", u.Email).Scan(ctx)
		if err == nil { // User exists
			log.Printf("User %s already exists, skipping.\n", u.Email)
			continue
		}
		if !errors.Is(err, sql.ErrNoRows) { // Some other error
			return err
		}
		// User does not exist, insert
		if _, err := db.NewInsert().Model(&u).Exec(ctx); err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, err)
		}
		log.Printf("Seeded user: %s\n", u.Email)
	}
	return nil
}

// Similar function for SeedURLs, ensuring user_id references existing seeded users.
