package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/config"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/internal/auth"
	"github.com/rigofekete/vhs-club-mvc/model"
	"github.com/rigofekete/vhs-club-mvc/repository"
)

type UserService interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	CreateUserBatch(context.Context, []*model.User) ([]*model.User, *int32, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UserLogin(ctx context.Context, user *model.User) (*model.User, error)
	GetAllUsers(context.Context) ([]*model.User, error)
	DeleteAllUsers(context.Context) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		repo: r,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	user.HashedPassword = hashedPassword
	return s.repo.Save(ctx, user)
}

func (s *userService) CreateUserBatch(ctx context.Context, users []*model.User) ([]*model.User, *int32, error) {
	for _, user := range users {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return nil, nil, err
		}
		user.Password = ""
		user.HashedPassword = hashedPassword
	}
	return s.repo.SaveBatch(ctx, users)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByPublicID(ctx, idUUID)
}

func (s *userService) UserLogin(ctx context.Context, user *model.User) (*model.User, error) {
	foundUser, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	valid, err := auth.CheckPasswordHash(user.Password, foundUser.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, apperror.ErrUserInvalidPW
	}

	token, err := auth.MakeJWT(foundUser.PublicID, foundUser.Role, config.AppConfig.JWTSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	foundUser.Token = token

	return foundUser, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) DeleteAllUsers(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}
