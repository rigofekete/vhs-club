package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/config"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/internal/database"
	"github.com/rigofekete/vhs-club-mvc/model"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) (*model.User, error)
	SaveBatch(ctx context.Context, users []*model.User) ([]*model.User, *int32, error)
	GetByID(ctx context.Context, id int32) (*model.User, error)
	GetByPublicID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	DeleteAll(ctx context.Context) error
}

type userRepository struct {
	DB *database.Queries
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		DB: config.AppConfig.DB,
		db: config.AppConfig.SQLDB,
	}
}

func (r *userRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	userParams := database.CreateUserParams{
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}

	dbUser, err := r.DB.CreateUser(ctx, userParams)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, apperror.ErrUserExists
		} else {
			return nil, err
		}
	}

	createdUser := &model.User{
		ID:        dbUser.ID,
		PublicID:  dbUser.PublicID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}
	return createdUser, nil
}

func (r *userRepository) SaveBatch(ctx context.Context, users []*model.User) ([]*model.User, *int32, error) {
	createdUsers := make([]*model.User, 0, len(users))
	existingCount := int32(0)
	for _, user := range users {
		userParams := database.CreateUserParams{
			Username:       user.Username,
			Email:          user.Email,
			HashedPassword: user.HashedPassword,
		}

		dbUser, err := r.DB.CreateUser(ctx, userParams)
		if err != nil {
			if isUniqueConstraintError(err) {
				existingCount++
				continue
			} else {
				return nil, nil, err
			}
		}

		createdUser := &model.User{
			ID:        dbUser.ID,
			PublicID:  dbUser.PublicID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
		}
		createdUsers = append(createdUsers, createdUser)
	}

	return createdUsers, &existingCount, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*model.User, error) {
	dbUser, err := r.DB.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperror.ErrUserNotFound, err)
	}
	user := &model.User{
		ID:        dbUser.ID,
		PublicID:  dbUser.PublicID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}
	return user, nil
}

func (r *userRepository) GetByPublicID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	dbUser, err := r.DB.GetUserByPublicID(ctx, id)
	if err != nil {
		return nil, apperror.ErrUserNotFound
	}

	user := &model.User{
		ID:        dbUser.ID,
		PublicID:  dbUser.PublicID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	dbUser, err := r.DB.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	user := &model.User{
		ID:             dbUser.ID,
		PublicID:       dbUser.PublicID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		Role:           dbUser.Role,
		HashedPassword: dbUser.HashedPassword,
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	dbUsers, err := r.DB.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*model.User, 0)
	for _, user := range dbUsers {
		u := &model.User{
			ID:        user.ID,
			PublicID:  user.PublicID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) DeleteAll(ctx context.Context) error {
	err := r.DB.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	return nil
}
