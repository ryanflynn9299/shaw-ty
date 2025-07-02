package models

import (
	"context"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel      `bun:"table:users alias:usr"`
	UUID               int64  `bun:"type:bigint,autoincrement:1"`
	FirstName          string `bun:"type:varchar(255)"`
	LastName           string `bun:"type:varchar(255)"`
	Email              string `bun:"type:varchar(255)"`
	Password           string `bun:"type:varchar(255)"`
	Salt               []byte `bun:"type:varchar(255)"`
	DateJoined         int64  `bun:"type:bigint"`
	DateModified       int64  `bun:"type:bigint"`
	LastPasswordUpdate int64  `bun:"type:bigint"`
	IsActive           bool   `bun:"type:bool"`
}

func (*User) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	_, err := query.DB().
		NewCreateIndex().
		Model((*User)(nil)).
		Index("email_idx").
		Column("email").
		Exec(ctx)
	return err
}
