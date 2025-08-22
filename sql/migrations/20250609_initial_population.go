package migrations

import (
	"URLShortener/internal/auth"
	"URLShortener/internal/storage/models"
	"context"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"time"
)

var Migrations = migrate.NewMigrations()
var _ bun.AfterCreateTableHook = (*models.User)(nil)

func init() {
	// Register the migration. The name should be unique and descriptive.
	// Bun uses the timestamp prefix to order migrations.
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// This is the UP migration.
		// It creates the tables and seeds the admin user.

		// 1. Create the 'users' table.
		_, err := db.NewCreateTable().
			Model((*models.User)(nil)).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}

		// 2. Create the 'short_links' table with a foreign key to 'users'.
		_, err = db.NewCreateTable().
			Model((*models.ShortLink)(nil)).
			IfNotExists().
			ForeignKey(`(creator_id) REFERENCES users (uuid) ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// 3. Create and insert the initial admin user.
		// IMPORTANT: Replace "supersecretpassword" with a secure password,
		// preferably loaded from an environment variable.
		password := "supersecretpassword42"
		salt, hashedPassword := auth.GetSaltAndHashedPassword(password)

		now := time.Now().Unix()
		adminUser := &models.User{
			UUID:               100,
			FirstName:          "Admin",
			LastName:           "User",
			Email:              "admin@example.com",
			Username:           "admin100",
			Salt:               salt,
			Password:           string(hashedPassword),
			DateJoined:         now,
			DateModified:       now,
			LastPasswordUpdate: now,
			IsActive:           true,
		}

		// We use an IGNORE clause to prevent errors if you run the migration
		// multiple times. The unique constraint on the email will prevent duplicates.
		_, err = db.NewInsert().
			Model(adminUser).
			Ignore().
			Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		// This is the DOWN migration.
		// It drops the tables in the reverse order of creation to respect foreign keys.

		// 1. Drop the 'short_links' table.
		_, err := db.NewDropTable().
			Model((*models.ShortLink)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			return err
		}

		// 2. Drop the 'users' table.
		_, err = db.NewDropTable().
			Model((*models.User)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
