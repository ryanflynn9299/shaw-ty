package db

import (
	"URLShortener/internal/storage/models"
	"context"
	"github.com/uptrace/bun"
)

// The Schema for the User Table is
// UUID	FIRSTNAME	LASTNAME	DATE_JOINED		EMAIL

// UserRepository defines the functionality for this repository
type UserRepository interface {
	Create(ctx context.Context, link *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	UpdateById(ctx context.Context, updatedUser *models.User) error
	SetUserActiveStatus(ctx context.Context, id int64, isActive bool) error
	DeleteById(ctx context.Context, id string) error
}

type UserRepositoryDB struct {
	db *bun.DB
}

func NewUserRepositoryDB(db *bun.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

func (r *UserRepositoryDB) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepositoryDB) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx, user)
	return user, err
}

func (r *UserRepositoryDB) GetByUsername(ctx context.Context, userName string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(&user).Where("username = ?", userName).Scan(ctx, user)
	return user, err
}

func (r *UserRepositoryDB) UpdateById(ctx context.Context, updatedUser *models.User) error {
	// TODO: verify this
	_, err := r.db.NewUpdate().Model(&updatedUser).Exec(ctx)
	return err
}

func (r *UserRepositoryDB) SetUserActiveStatus(ctx context.Context, id int64, isActive bool) error {
	user := new(models.User)
	_, err := r.db.NewUpdate().Model(&user).Where("id = ?", id).Set("is_active", isActive).Exec(ctx)
	return err
}

func (r *UserRepositoryDB) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.NewSelect().Model(&users).Scan(ctx, users)
	return users, err
}

func (r *UserRepositoryDB) DeleteById(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Where("uuid = ?", id).Exec(ctx)
	return err
}
