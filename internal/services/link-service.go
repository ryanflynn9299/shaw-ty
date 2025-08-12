package services

import (
	"URLShortener/internal/core/encoder"
	idgenerator "URLShortener/internal/core/id-generator"
	"URLShortener/internal/storage/db"
	"URLShortener/internal/storage/models"
	"context"
	"fmt"
	"time"
)

type LinkServiceIfc interface {
	CreateLink(url string, userId int64, customLink string, expiresAfter int, sfGen *idgenerator.SnowflakeGenerator) (string, error)
	GetLinkById(linkId string) (*models.ShortLink, error)
	GetAllLinks() ([]models.ShortLink, error)
	UpdateLinkById(id int, shortLink string, fullUrl string, isActive bool) (int, error)
	DeactivateLink(linkId string) error
	DeleteLink(linkId string) error
}

type LinkService struct {
	repo    db.ShortLinkRepository
	context context.Context
}

func (l LinkService) CreateLink(url string, userId int64, customLink string, expiresAfter int, sfGen *idgenerator.SnowflakeGenerator) (string, error) {
	linkId := sfGen.NextId()
	createDate := time.Now()
	var expiresDate int64
	if expiresAfter != 0 {
		expiresDate = time.Now().AddDate(expiresAfter, 0, 0).UnixMilli()
	} else {
		expiresDate = time.Now().AddDate(50, 0, 1).UnixMilli()
	}
	shortName := encoder.Base63Encode(int64(linkId))

	newShortLink := &models.ShortLink{
		ID:             int64(linkId),
		CreatedDate:    createDate.UnixMilli(),
		ExpirationDate: expiresDate,
		FullURL:        url,
		ShortenedCode:  shortName,
		CustomCode:     customLink,
		CreatorId:      userId,
		IsActive:       true,
	}

	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	err := l.repo.Create(ctx, newShortLink)
	if err != nil {
		return "", err
	}

	return shortName, nil
}

func (l LinkService) GetLinkById(linkId string) (*models.ShortLink, error) {
	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	id, err := l.repo.GetByID(ctx, linkId)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (l LinkService) GetLinkByCode(code string) (*models.ShortLink, error) {
	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	id, err := l.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (l LinkService) GetAllLinks() ([]models.ShortLink, error) {
	ctx, cancel := context.WithTimeout(l.context, 10*time.Second)
	defer cancel()

	ids, err := l.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (l LinkService) UpdateLink(id int, shortLink string, fullUrl string, isActive bool) (int, error) {
	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	link, err := l.repo.GetByID(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		panic(err) // TODO: don't panic, handle gracefully
	}

	// Check for any changes, returning success for idempotency early if no changes were made.
	// 		No special reasoning is given for success to avoid revealing DB state to users
	var didAnyDataChange = link.CustomCode != shortLink || link.FullURL != fullUrl || link.IsActive != isActive
	if !didAnyDataChange {
		// Nothing changed, return success
		return 0, nil
	}

	if shortLink != "" {
		link.ShortenedCode = shortLink
	}
	if fullUrl != "" {
		link.FullURL = fullUrl
	}
	link.IsActive = isActive
	link.DateModified = time.Now().UnixMilli()

	err = l.repo.UpdateById(ctx, link)
	if err != nil {
		return 1, err
	}

	// link data successfully updated, return success
	return 0, nil
}

func (l LinkService) DeactivateLink(linkId string) error {
	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	err := l.repo.SetActiveById(ctx, linkId, false)
	return err
}

func (l LinkService) DeleteLink(linkId string) error {
	ctx, cancel := context.WithTimeout(l.context, 1*time.Second)
	defer cancel()

	err := l.repo.DeleteById(ctx, linkId)
	return err
}

func NewLinkService(repo db.ShortLinkRepository, ctx context.Context) LinkService {
	return LinkService{repo, ctx}
}
