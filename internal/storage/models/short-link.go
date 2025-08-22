package models

import "github.com/uptrace/bun"

type ShortLink struct {
	bun.BaseModel  `bun:"table:short_links"`
	ID             int64  `bun:",pk,type:bigint"` // not autoincrement because of snowflake generator
	CreatedDate    int64  `bun:"created_date,type:bigint"`
	DateModified   int64  `bun:"date_modified,type:bigint"`
	ExpirationDate int64  `bun:"expiration_date,type:bigint"`
	FullURL        string `bun:"full_url,type:text"`
	ShortenedCode  string `bun:"shortened_code,type:text"`       // base63 encoded string
	CreatorId      int64  `bun:"creator_id,type:bigint,notnull"` // the uuid of the user who created the shortlink
	CustomCode     string `bun:"custom_code,type:text"`          // a user-requested string
	IsActive       bool   `bun:"is_active,type:boolean"`

	Creator *User `bun:"rel:belongs-to,join:creator_id=uuid"`
}
