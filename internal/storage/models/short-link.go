package models

import "github.com/uptrace/bun"

type ShortLink struct {
	bun.BaseModel  `bun:"table:short_links,alias:short_link"`
	ID             int64  `bun:"id,type:bigint"` // not autoincrement because of snowflake generator
	CreatedDate    int64  `bun:"created_date,type:bigint"`
	DateModified   int64  `bun:"date_modified,type:bigint"`
	ExpirationDate int64  `bun:"expiration_date,type:bigint"`
	FullURL        string `bun:"full_url,type:text"`
	ShortenedCode  string `bun:"shortened_code,type:text"` // base63 encoded string
	CreatorId      int64  `bun:"creator_id,type:bigint"`   // the uuid of the user who created the shortlink
	CustomCode     string `bun:"custom_code,type:text"`    // a user-requested string
	IsActive       bool   `bun:"is_active,type:boolean"`
}
