package models

import (
	"context"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel      `bun:"table:users"`
	UUID               int64  `bun:"uuid,pk,autoincrement"`
	FirstName          string `bun:"type:varchar(255)"`
	LastName           string `bun:"type:varchar(255)"`
	Email              string `bun:"type:varchar(255),unique,notnull"`
	Username           string `bun:"type:varchar(255)"`
	Password           string `bun:"type:varchar(255),notnull"`
	Salt               string `bun:"type:varchar(255)"`
	DateJoined         int64  `bun:"type:bigint"`
	DateModified       int64  `bun:"type:bigint"`
	LastPasswordUpdate int64  `bun:"type:bigint"`
	IsActive           bool   `bun:"type:bool"`

	ShortLinks []*ShortLink `bun:"rel:has-many,join:uuid=creator_id"`
}

type SafeUser struct {
	bun.BaseModel      `bun:"table:users"`
	UUID               int64  `bun:"column:pk,type:bigint,autoincrement"`
	FirstName          string `bun:"type:varchar(255)"`
	LastName           string `bun:"type:varchar(255)"`
	Email              string `bun:"type:varchar(255)"`
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
