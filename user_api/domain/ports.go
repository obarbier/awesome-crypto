package domain

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=./mock_user_repo.go -package=domain github.com/obarbier/awesome-crypto/user_api/domain UserRepository
type UserRepository interface {
	Save(ctx context.Context, firstName, lastName, userId, passwordHash string) (User, error)
	Update(ctx context.Context, id, firstName, lastName, userId, passwordHash string) (User, error)
	Get(ctx context.Context, userId string) (User, error)
	Delete(ctx context.Context, userId string) error
}

type IUserService interface {
	CreateUser(ctx context.Context, firstName, lastName, userId, password string) (User, error)
	UpdateUser(ctx context.Context, id, firstName, lastName, userId, password string) (User, error)
	GetUserById(ctx context.Context, userId string) (User, error)
	DeleteUser(ctx context.Context, userId string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, firstName, lastName, userId, password string) (User, error) {
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return User{}, hashedPasswordErr
	}
	hashVerifyErr := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if hashVerifyErr != nil {
		return User{}, hashVerifyErr
	}
	return u.repo.Save(ctx, firstName, lastName, userId, string(hashedPassword))
}

func (u *UserService) UpdateUser(ctx context.Context, id, firstName, lastName, userId, password string) (User, error) {
	return u.repo.Update(ctx, id, firstName, lastName, userId, password)
}

func (u *UserService) GetUserById(ctx context.Context, userId string) (User, error) {
	return u.repo.Get(ctx, userId)
}

func (u *UserService) DeleteUser(ctx context.Context, userId string) error {
	return u.repo.Delete(ctx, userId)
}
