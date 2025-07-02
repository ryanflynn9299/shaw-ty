package migrations

import (
	"URLShortener/internal/storage/models"
	"context"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()
var _ bun.AfterCreateTableHook = (*models.User)(nil)

func init() {

	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			// up migration
			_, err := db.NewCreateTable().
				Model((*models.User)(nil)).
				IfNotExists().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = db.NewCreateTable().
				Model((*models.ShortLink)(nil)).
				IfNotExists().
				ForeignKey("(creator_id) REFERENCES user (uuid) ON DELETE CASCADE").
				Exec(ctx)

			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			// down migration
			_, err := db.NewDropTable().Model((*models.ShortLink)(nil)).IfExists().Exec(ctx)
			if err != nil {
				return err
			}

			_, err = db.NewDropTable().Model((*models.User)(nil)).IfExists().Exec(ctx)
			return err
		})
}
