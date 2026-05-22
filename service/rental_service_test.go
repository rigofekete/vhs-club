package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/model"
	"github.com/rigofekete/vhs-club-mvc/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRentalRepository struct {
	mock.Mock
}

func NewRentalMockRepository() *mockRentalRepository {
	return &mockRentalRepository{}
}

func (m *mockRentalRepository) Save(ctx context.Context, userID, tapeID int32) (*model.Rental, error) {
	args := m.Called(ctx, userID, tapeID)
	if r := args.Get(0); r != nil {
		return r.(*model.Rental), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRentalRepository) ReturnTape(ctx context.Context, rentalID uuid.UUID, userID int32) error {
	args := m.Called(ctx, rentalID, userID)
	return args.Error(0)
}

func (m *mockRentalRepository) GetUserActive(ctx context.Context, userID int32) ([]*model.Rental, error) {
	args := m.Called(ctx, userID)
	if rentals := args.Get(0); rentals != nil {
		return rentals.([]*model.Rental), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRentalRepository) GetAllActive(ctx context.Context) ([]*model.Rental, error) {
	args := m.Called(ctx)
	if rentals := args.Get(0); rentals != nil {
		return rentals.([]*model.Rental), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRentalRepository) GetActiveRentCountByTape(ctx context.Context, id int32) (*int64, error) {
	args := m.Called(ctx, id)
	if count := args.Get(0); count != nil {
		return count.(*int64), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRentalRepository) GetActiveRentCountByUser(ctx context.Context, id int32) (*int64, error) {
	args := m.Called(ctx, id)
	if rental := args.Get(0); rental != nil {
		return rental.(*int64), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRentalRepository) DeleteAllRentals(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func Test_RentTape_Success(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	tapeUUID := uuid.New()
	userUUID := uuid.New()
	tapeID := int32(8)
	userID := int32(14)
	userRentCount := int64(1)
	tapeRentCount := int64(0)

	user := &model.User{
		ID:       userID,
		Username: "NoamChomsky",
		Email:    "LAD@mit.edu",
	}

	tape := &model.Tape{
		ID:       tapeID,
		Title:    "The Matrix",
		Director: "Wachowski sisters",
		Genre:    "Cyberpunk",
		Quantity: 1,
	}

	rentalID := int32(23)
	dbRental := &model.Rental{
		ID:        rentalID,
		TapeTitle: "The Matrix",
		Username:  "NoamChomsky",
	}

	ctx := context.Background()
	mockTapeRepo.On("GetByPublicID", ctx, tapeUUID).Return(tape, nil)
	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(user, nil)
	mockRentalRepo.On("GetActiveRentCountByUser", ctx, userID).Return(&userRentCount, nil)
	mockRentalRepo.On("GetActiveRentCountByTape", ctx, tapeID).Return(&tapeRentCount, nil)
	mockRentalRepo.On("Save", ctx, tapeID, userID).Return(dbRental, nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rental, err := svc.RentTape(ctx, tapeUUID.String(), userUUID.String())

	assert.Nil(t, err)
	assert.Equal(t, dbRental, rental)

	mockRentalRepo.AssertExpectations(t)
}

func Test_RentTape_TapeNotFound(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockUserRepo := NewUserMockRepository()
	mockTapeRepo := NewTapeMockRepository()

	userUUID := uuid.New()
	tapeUUID := uuid.New()

	ctx := context.Background()
	mockTapeRepo.On("GetByPublicID", ctx, tapeUUID).Return(nil, apperror.ErrTapeNotFound)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rental, err := svc.RentTape(ctx, tapeUUID.String(), userUUID.String())

	assert.Error(t, err)
	assert.Equal(t, err, apperror.ErrTapeNotFound)
	assert.Nil(t, rental)

	mockRentalRepo.AssertExpectations(t)
}

func Test_RentTape_Fail_UserNotFound(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockUserRepo := NewUserMockRepository()
	mockTapeRepo := NewTapeMockRepository()

	userUUID := uuid.New()
	tapeUUID := uuid.New()

	ctx := context.Background()
	mockTapeRepo.On("GetByPublicID", ctx, tapeUUID).Return(&model.Tape{}, nil)
	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(nil, apperror.ErrUserNotFound)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rental, err := svc.RentTape(ctx, tapeUUID.String(), userUUID.String())

	assert.Error(t, err)
	assert.Equal(t, err, apperror.ErrUserNotFound)
	assert.Nil(t, rental)

	mockRentalRepo.AssertExpectations(t)
}

func Test_RentTape_Fail_TapeUnavailable(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockUserRepo := NewUserMockRepository()
	mockTapeRepo := NewTapeMockRepository()

	userUUID := uuid.New()
	tapeUUID := uuid.New()
	tapeID := int32(14)
	tapeRentCount := int64(1)

	returnedTape := &model.Tape{
		ID:       tapeID,
		Quantity: 1,
	}

	ctx := context.Background()
	mockTapeRepo.On("GetByPublicID", ctx, tapeUUID).Return(returnedTape, nil)
	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(&model.User{}, nil)
	mockRentalRepo.On("GetActiveRentCountByTape", ctx, tapeID).Return(&tapeRentCount, nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rental, err := svc.RentTape(ctx, tapeUUID.String(), userUUID.String())

	assert.Error(t, err)
	assert.Equal(t, err, apperror.ErrTapeUnavailable)
	assert.Nil(t, rental)

	mockRentalRepo.AssertExpectations(t)
}

func Test_RentTape_Fail_MaxRentalsPerUser(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockUserRepo := NewUserMockRepository()
	mockTapeRepo := NewTapeMockRepository()

	userUUID := uuid.New()
	tapeUUID := uuid.New()
	tapeID := int32(14)
	userID := int32(80)
	tapeRentCount := int64(1)
	// User set to have max allowed active rented tapes (2)
	userRentCount := int64(2)

	returnedTape := &model.Tape{
		ID:       tapeID,
		Quantity: 2,
	}

	returnedUser := &model.User{
		ID: userID,
	}

	ctx := context.Background()
	mockTapeRepo.On("GetByPublicID", ctx, tapeUUID).Return(returnedTape, nil)
	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(returnedUser, nil)
	mockRentalRepo.On("GetActiveRentCountByTape", ctx, tapeID).Return(&tapeRentCount, nil)
	mockRentalRepo.On("GetActiveRentCountByUser", ctx, userID).Return(&userRentCount, nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rental, err := svc.RentTape(ctx, tapeUUID.String(), userUUID.String())

	assert.Error(t, err)
	assert.Equal(t, err, apperror.ErrMaxRentalsPerUser)
	assert.Nil(t, rental)

	mockRentalRepo.AssertExpectations(t)
}

func Test_ReturnTape_Success(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	ctx := context.Background()
	rentalUUID := uuid.New()
	userUUID := uuid.New()
	userID := int32(10)
	user := &model.User{
		ID: userID,
	}
	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(user, nil)
	mockRentalRepo.On("ReturnTape", ctx, rentalUUID, user.ID).Return(nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	err := svc.ReturnTape(ctx, userUUID.String(), rentalUUID.String())

	assert.Nil(t, err)

	mockRentalRepo.AssertExpectations(t)
}

func Test_ReturnTape_UserNotFound(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	ctx := context.Background()
	rentalUUID := uuid.New()
	userUUID := uuid.New()

	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(nil, apperror.ErrUserNotFound)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	err := svc.ReturnTape(ctx, userUUID.String(), rentalUUID.String())

	assert.Error(t, err)
	assert.Error(t, err, apperror.ErrUserNotFound)

	mockRentalRepo.AssertExpectations(t)
}

func Test_GetUserActiveRentals(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	dbRentals := []*model.Rental{
		{
			Username: "RonGilbert",
		},
		{
			Username: "JohnCarmack",
		},
	}

	ctx := context.Background()
	userID := int32(32)
	userUUID := uuid.New()
	user := &model.User{
		ID:       userID,
		PublicID: userUUID,
	}

	mockUserRepo.On("GetByPublicID", ctx, userUUID).Return(user, nil)
	mockRentalRepo.On("GetUserActive", ctx, user.ID).Return(dbRentals, nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rentals, err := svc.GetUserActiveRentals(ctx, userUUID.String())

	assert.Nil(t, err)
	assert.Equal(t, dbRentals, rentals)

	mockRentalRepo.AssertExpectations(t)
}

func Test_GetAllActiveRentals(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	dbRentals := []*model.Rental{
		{
			Username: "RonGilbert",
		},
		{
			Username: "JohnCarmack",
		},
	}

	ctx := context.Background()
	mockRentalRepo.On("GetAllActive", ctx).Return(dbRentals, nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	rentals, err := svc.GetAllActiveRentals(ctx)

	assert.Nil(t, err)
	assert.Equal(t, dbRentals, rentals)

	mockRentalRepo.AssertExpectations(t)
}

func Test_DeleteAllRentals(t *testing.T) {
	mockRentalRepo := NewRentalMockRepository()
	mockTapeRepo := NewTapeMockRepository()
	mockUserRepo := NewUserMockRepository()

	ctx := context.Background()
	mockRentalRepo.On("DeleteAllRentals", ctx).Return(nil)

	svc := service.NewRentalService(mockRentalRepo, mockTapeRepo, mockUserRepo)
	err := svc.DeleteAllRentals(ctx)

	assert.Nil(t, err)

	mockRentalRepo.AssertExpectations(t)
}
