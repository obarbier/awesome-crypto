package domain

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=./mock_user_repo.go -package=domain github.com/obarbier/awesome-crypto/user_api/domain UserRepository
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	UpdateByID(ctx context.Context, id string, user *User) error
	FindById(ctx context.Context, id string) (*User, error)
	DeleteById(ctx context.Context, id string) error
}

type IUserService interface {
	CreateUser(ctx context.Context, firstName, lastName, userId, password string) (*User, error)
	UpdateUser(ctx context.Context, id, firstName, lastName, userId, password string) error
	GetUserById(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, firstName, lastName, userId, password string) (*User, error) {
	hashedPassword, encryptErr := encryptPassword(password)
	if encryptErr != nil {
		return nil, encryptErr
	}
	id, uuidErr := uuid.NewUUID()
	if uuidErr != nil {
		return nil, uuidErr
	}
	user := &User{
		Id:           id.String(),
		FirstName:    firstName,
		LastName:     lastName,
		UserId:       userId,
		PasswordHash: string(hashedPassword),
	}

	savedErr := u.repo.Save(ctx, user)
	if savedErr != nil {
		return nil, savedErr
	}
	return user, nil

}

func (u *UserService) UpdateUser(ctx context.Context, id, firstName, lastName, userId, password string) error {
	updates := &User{}
	if firstName != "" {
		updates.FirstName = firstName
	}
	if lastName != "" {
		updates.LastName = lastName
	}
	if userId != "" {
		updates.UserId = userId
	}
	if password != "" {
		hashedPassword, encryptErr := encryptPassword(password)
		if encryptErr != nil {
			return encryptErr
		}
		updates.PasswordHash = string(hashedPassword)
	}
	return u.repo.UpdateByID(ctx, id, updates)
}

func (u *UserService) GetUserById(ctx context.Context, id string) (*User, error) {
	return u.repo.FindById(ctx, id)
}

func (u *UserService) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteById(ctx, id)
}
