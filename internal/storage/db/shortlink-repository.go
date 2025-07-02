package db

import (
	"URLShortener/internal/storage/models"
	"context"
	"github.com/uptrace/bun"
)

// The schema for the main data table is:
// ID	CREATION_DATE	EXPIRATION_DATE	FULL_URL	SHORTENED_CODE	CREATOR_ID	CUSTOM_CODE	IS_ACTIVE

// ShortLinkRepository Define the functionality of this repository
type ShortLinkRepository interface {
	Create(ctx context.Context, link *models.ShortLink) error
	GetByID(ctx context.Context, id string) (*models.ShortLink, error)
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetAll(ctx context.Context) ([]models.ShortLink, error)
	SetActiveById(ctx context.Context, id string, isActive bool) error
	DeleteById(ctx context.Context, id string) error
	// update link - set all fields changed
}

type ShortLinkRepositoryDB struct {
	db *bun.DB
}

func NewShortLinkRepositoryDB(db *bun.DB) *ShortLinkRepositoryDB {
	return &ShortLinkRepositoryDB{db: db}
}

func (r *ShortLinkRepositoryDB) Create(ctx context.Context, link *models.ShortLink) error {
	_, err := r.db.NewInsert().Model(link).Exec(ctx)
	return err
}

func (r *ShortLinkRepositoryDB) GetByID(ctx context.Context, id string) (*models.ShortLink, error) {
	shortlink := new(models.ShortLink)
	err := r.db.NewSelect().Model(shortlink).Where("id = ?", id).Scan(ctx, shortlink)
	return shortlink, err
}

func (r *ShortLinkRepositoryDB) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {
	shortlink := new(models.ShortLink)
	err := r.db.NewSelect().Model(shortlink).
		Where("shortened_code = ?", code).WhereOr("custom_code = ?", code).
		Scan(ctx, shortlink)
	return shortlink, err
}

func (r *ShortLinkRepositoryDB) GetAll(ctx context.Context) ([]models.ShortLink, error) {
	var shortlinks []models.ShortLink
	err := r.db.NewSelect().Model(&shortlinks).Scan(ctx, shortlinks)
	return shortlinks, err
}

func (r *ShortLinkRepositoryDB) SetActiveById(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.NewUpdate().SetColumn("isActive", "?", isActive).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *ShortLinkRepositoryDB) DeleteById(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Where("id = ?", id).Exec(ctx)
	return err
}
