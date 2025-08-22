package services

import (
	"URLShortener/internal/storage/db"
	"URLShortener/internal/storage/models"
	"context"
	"fmt"
	"log"
	"time"
)

type UserService interface {
	RegisterUser(first string, last string, email string, username string, password string, salt string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(id int) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	UpdateUserById(id int, username string, password string, first string, last string, email string, isActive bool) (int, error)
	UpdateCompleteUserById(id int, username string, password string, salt string, UUID int64, first string, last string, email string, dateJoined int64, dateModified int64, lastPasswordUpdate int64, isActive bool) (int, error)
	DeactivateUserById(id int) (bool, error)
	ReactivateUserById(id int) (bool, error)
	DeleteUserById(id int) error
}

type userService struct {
	repo    db.UserRepository
	context context.Context
}

func NewUserService(repo db.UserRepository, context context.Context) UserService {
	return &userService{
		repo:    repo,
		context: context,
	}
}

// RegisterUser invokes the Create operation on the user table, passing the provided data in a model
func (s *userService) RegisterUser(first string, last string, email string, username string, hashedPassword string, salt string) (*models.User, error) {
	newUser := &models.User{
		FirstName:          first,
		LastName:           last,
		Email:              email,
		Username:           username,
		Password:           hashedPassword,
		Salt:               salt,
		LastPasswordUpdate: time.Now().UnixMilli(),
		DateJoined:         time.Now().UnixMilli(),
		DateModified:       time.Now().UnixMilli(),
		IsActive:           true,
	}

	ctx, cancel := context.WithTimeout(s.context, 10*time.Second)
	defer cancel()
	err := s.repo.Create(ctx, newUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newUser, nil
}

func (s *userService) FindUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(s.context, 10*time.Second)
	defer cancel()

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(s.context, 10*time.Second)
	defer cancel()

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

// GetUserById retrieves the user with given id or returns an error
func (s *userService) GetUserById(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(s.context, 1*time.Second)
	defer cancel()

	user, err := s.repo.GetByID(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserById(id int, username string, password string, first string, last string, email string, isActive bool) (int, error) {
	ctx, cancel := context.WithTimeout(s.context, 1*time.Second)
	defer cancel()

	// TODO: enforce RBAP rules

	// Retrieve existing entry for comparison
	user, err := s.repo.GetByID(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		panic(err) // TODO: don't panic, handle gracefully
	}

	// Check for any changes, returning success for idempotency early if no changes were made.
	// 		No special reasoning is given for success to avoid revealing DB state to users
	var didPwChange = user.Password != password && password != ""
	var didAnyChange = user.Email != email || user.FirstName != first || user.LastName != last || didPwChange
	if !didAnyChange {
		// do nothing, return success
		return 0, nil
	}

	if didPwChange {
		user.LastPasswordUpdate = time.Now().UnixMilli()
	}

	// Apply changes to user object
	if didPwChange {
		user.Password = password
	}
	if username != "" {
		user.Username = username
	}
	if first != "" {
		user.FirstName = first
	}
	if last != "" {
		user.LastName = last
	}
	if email != "" {
		user.Email = email
	}
	if isActive != user.IsActive {
		user.IsActive = isActive
	}

	user.DateModified = time.Now().UnixMilli()

	err = s.repo.UpdateById(ctx, user)
	if err != nil {
		return 1, err
	}

	// user modified, return success
	return 0, nil
}

func (s *userService) UpdateCompleteUserById(id int, username string, password string, salt string, UUID int64, first string, last string, email string, dateJoined int64, dateModified int64, lastPasswordUpdate int64, isActive bool) (int, error) {
	ctx, cancel := context.WithTimeout(s.context, 1*time.Second)
	defer cancel()

	// TODO: enforce RBAP rules

	// Retrieve existing entry for comparison
	user, err := s.repo.GetByID(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		panic(err) // TODO: don't panic, handle gracefully
	}

	// Validate that the calling user has RBAP over changee id

	// TODO: fix this didAnyChange logic
	var didPwChange = user.Password != password && password != ""
	var didAnyChange = user.Email != email || user.FirstName != first || user.LastName != last || didPwChange
	if !didAnyChange {
		// do nothing, return success
		return 0, nil
	}

	if didPwChange {
		user.LastPasswordUpdate = time.Now().UnixMilli()
	}

	// Apply changes to user object
	if didPwChange {
		user.Password = password
	}
	if username != "" {
		user.Username = username
	}
	if UUID != user.UUID {
		user.UUID = UUID
	}
	if salt != user.Salt {
		user.Salt = salt
	}
	if first != "" {
		user.FirstName = first
	}
	if last != "" {
		user.LastName = last
	}
	if email != "" {
		user.Email = email
	}
	if dateJoined != user.DateJoined && dateJoined != 0 {
		user.DateJoined = dateJoined
	}
	if dateModified != user.DateModified && dateModified != 0 {
		user.DateModified = dateModified
	}
	if lastPasswordUpdate != user.LastPasswordUpdate && lastPasswordUpdate != 0 {
		user.DateModified = dateModified
	}
	if isActive != user.IsActive {
		user.IsActive = isActive
	}

	err = s.repo.UpdateById(ctx, user)
	if err != nil {
		return 1, err
	}

	// user modified, return success
	return 0, nil
}

func (s *userService) DeactivateUserById(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(s.context, 10*time.Second)
	defer cancel()

	// TODO: enforce RBAP rules

	err := s.repo.SetUserActiveStatus(ctx, int64(id), false)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) ReactivateUserById(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(s.context, 10*time.Second)
	defer cancel()

	// TODO: enforce RBAP rules

	// TODO: prompt new password dialog? -- or let business/application logic handle
	// TODO: expire current password, application must read expiry and show dialog
	err := s.repo.SetUserActiveStatus(ctx, int64(id), true)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteUserById deletes the user specified by ID and errors if no such user exists
func (s *userService) DeleteUserById(id int) error {
	ctx, cancel := context.WithTimeout(s.context, 1*time.Second)
	defer cancel()

	// TODO: enforce RBAP rules

	err := s.repo.DeleteById(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
