package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/model"
	"github.com/rigofekete/vhs-club-mvc/repository"
)

type RentalService interface {
	RentTape(ctx context.Context, tapeID string, userID string) (*model.Rental, error)
	ReturnTape(ctx context.Context, userID, rentalID string) error
	GetUserActiveRentals(ctx context.Context, userID string) ([]*model.Rental, error)
	GetAllActiveRentals(ctx context.Context) ([]*model.Rental, error)
	DeleteAllRentals(ctx context.Context) error
}

type rentalService struct {
	tapeRepo   repository.TapeRepository
	userRepo   repository.UserRepository
	rentalRepo repository.RentalRepository
}

func NewRentalService(r repository.RentalRepository, t repository.TapeRepository, u repository.UserRepository) RentalService {
	return &rentalService{
		rentalRepo: r,
		tapeRepo:   t,
		userRepo:   u,
	}
}

// Business logic local constants
const maxRentalsPerUser = 2

func (s *rentalService) RentTape(ctx context.Context, tapePublicID, userPublicID string) (*model.Rental, error) {
	tapeUUID, err := uuid.Parse(tapePublicID)
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.Parse(userPublicID)
	if err != nil {
		return nil, err
	}

	tape, err := s.tapeRepo.GetByPublicID(ctx, tapeUUID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByPublicID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	countByTape, err := s.rentalRepo.GetActiveRentCountByTape(ctx, tape.ID)
	if err != nil {
		return nil, err
	}

	if int32(*countByTape) >= tape.Quantity {
		return nil, apperror.ErrTapeUnavailable
	}

	countByUser, err := s.rentalRepo.GetActiveRentCountByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if int32(*countByUser) >= maxRentalsPerUser {
		return nil, apperror.ErrMaxRentalsPerUser
	}

	return s.rentalRepo.Save(ctx, tape.ID, user.ID)
}

func (s *rentalService) ReturnTape(ctx context.Context, userPublicID, rentalPublicID string) error {
	rentalUUID, err := uuid.Parse(rentalPublicID)
	if err != nil {
		return err
	}

	userUUID, err := uuid.Parse(userPublicID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByPublicID(ctx, userUUID)
	if err != nil {
		return err
	}

	return s.rentalRepo.ReturnTape(ctx, rentalUUID, user.ID)
}

func (s *rentalService) GetUserActiveRentals(ctx context.Context, userPublicID string) ([]*model.Rental, error) {
	userUUID, err := uuid.Parse(userPublicID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByPublicID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return s.rentalRepo.GetUserActive(ctx, user.ID)
}

func (s *rentalService) GetAllActiveRentals(ctx context.Context) ([]*model.Rental, error) {
	return s.rentalRepo.GetAllActive(ctx)
}

func (s *rentalService) DeleteAllRentals(ctx context.Context) error {
	return s.rentalRepo.DeleteAllRentals(ctx)
}
