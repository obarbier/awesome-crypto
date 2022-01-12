package domain

import (
	"context"
	"fmt"
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

var _ IUserService = (*UserService)(nil)

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, firstName, lastName, userId, password string) (*User, error) {
	if len(firstName) == 0 {
		return nil, fmt.Errorf("FirstName cannot be empty")
	}
	if len(lastName) == 0 {
		return nil, fmt.Errorf("LastName cannot be empty")
	}
	if len(userId) == 0 {
		return nil, fmt.Errorf("UserId cannot be empty")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}
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
