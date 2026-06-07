package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/rigofekete/vhs-club-mvc/config"
	"github.com/rigofekete/vhs-club-mvc/internal/apperror"
	"github.com/rigofekete/vhs-club-mvc/internal/database"
	"github.com/rigofekete/vhs-club-mvc/model"
)

type TapeRepository interface {
	Save(ctx context.Context, tape *model.Tape) (*model.Tape, error)
	SaveBatch(ctx context.Context, tapes []*model.Tape) ([]*model.Tape, *int32, error)
	GetAll(ctx context.Context) ([]*model.Tape, error)
	GetByID(ctx context.Context, id int32) (*model.Tape, error)
	GetByPublicID(ctx context.Context, id uuid.UUID) (*model.Tape, error)
	Update(ctx context.Context, updateTape *model.UpdateTape) (*model.Tape, error)
	Delete(ctx context.Context, id int32) error
	DeleteAll(ctx context.Context) error
}

type tapeRepository struct {
	DB *database.Queries
	db *sql.DB
}

func NewTapeRepository() TapeRepository {
	return &tapeRepository{
		DB: config.AppConfig.DB,
		db: config.AppConfig.SQLDB,
	}
}

func (r *tapeRepository) Save(ctx context.Context, tape *model.Tape) (*model.Tape, error) {
	tapeParams := database.CreateTapeParams{
		Title:    tape.Title,
		Director: tape.Director,
		Genre:    tape.Genre,
		Quantity: tape.Quantity,
	}

	dbTape, err := r.DB.CreateTape(ctx, tapeParams)
	if err != nil {
		return nil, err
	}

	savedTape := &model.Tape{
		ID:        dbTape.ID,
		PublicID:  dbTape.PublicID,
		CreatedAt: dbTape.CreatedAt,
		UpdatedAt: dbTape.UpdatedAt,
		Title:     dbTape.Title,
		Director:  dbTape.Director,
		Genre:     dbTape.Genre,
		Quantity:  dbTape.Quantity,
	}
	return savedTape, nil
}

func (r *tapeRepository) SaveBatch(ctx context.Context, tapes []*model.Tape) ([]*model.Tape, *int32, error) {
	createdTapes := make([]*model.Tape, 0, len(tapes))
	existingCount := int32(0)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	qtx := r.DB.WithTx(tx)
	for _, tape := range tapes {
		tapeParams := database.CreateTapeParams{
			Title:    tape.Title,
			Director: tape.Director,
			Genre:    tape.Genre,
			Quantity: tape.Quantity,
		}

		dbTape, err := qtx.CreateTape(ctx, tapeParams)
		if err != nil {
			if isUniqueConstraintError(err) {
				existingCount++
				continue
			} else {
				return nil, nil, err
			}
		}

		createdTape := &model.Tape{
			ID:        dbTape.ID,
			PublicID:  dbTape.PublicID,
			CreatedAt: dbTape.CreatedAt,
			UpdatedAt: dbTape.UpdatedAt,
			Title:     dbTape.Title,
			Director:  dbTape.Director,
			Genre:     dbTape.Genre,
			Quantity:  dbTape.Quantity,
		}
		createdTapes = append(createdTapes, createdTape)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return createdTapes, &existingCount, nil
}

func (r *tapeRepository) GetAll(ctx context.Context) ([]*model.Tape, error) {
	dbTapes, err := r.DB.GetTapes(ctx)
	if err != nil {
		return nil, err
	}
	tapes := make([]*model.Tape, 0)
	for _, tape := range dbTapes {
		t := &model.Tape{
			ID:        tape.ID,
			PublicID:  tape.PublicID,
			CreatedAt: tape.CreatedAt,
			UpdatedAt: tape.UpdatedAt,
			Title:     tape.Title,
			Director:  tape.Director,
			Genre:     tape.Genre,
			Quantity:  tape.Quantity,
		}
		tapes = append(tapes, t)
	}
	return tapes, nil
}

func (r *tapeRepository) GetByID(ctx context.Context, id int32) (*model.Tape, error) {
	dbTape, err := r.DB.GetTapeByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrTapeNotFound
	}

	tape := &model.Tape{
		ID:        dbTape.ID,
		PublicID:  dbTape.PublicID,
		CreatedAt: dbTape.CreatedAt,
		UpdatedAt: dbTape.UpdatedAt,
		Title:     dbTape.Title,
		Director:  dbTape.Director,
		Genre:     dbTape.Genre,
		Quantity:  dbTape.Quantity,
	}

	return tape, nil
}

func (r *tapeRepository) GetByPublicID(ctx context.Context, id uuid.UUID) (*model.Tape, error) {
	dbTape, err := r.DB.GetTapeFromPublicID(ctx, id)
	if err != nil {
		return nil, apperror.ErrTapeNotFound
	}
	tape := &model.Tape{
		ID:        dbTape.ID,
		PublicID:  dbTape.PublicID,
		CreatedAt: dbTape.CreatedAt,
		UpdatedAt: dbTape.UpdatedAt,
		Title:     dbTape.Title,
		Director:  dbTape.Director,
		Genre:     dbTape.Genre,
		Quantity:  dbTape.Quantity,
	}
	return tape, nil
}

func (r *tapeRepository) Update(ctx context.Context, updateTape *model.UpdateTape) (*model.Tape, error) {
	dbUpdateParams := database.UpdateTapeParams{
		ID:       updateTape.ID,
		Title:    toNullString(updateTape.Title),
		Director: toNullString(updateTape.Director),
		Genre:    toNullString(updateTape.Genre),
		Quantity: toNullInt32(updateTape.Quantity),
	}

	dbTape, err := r.DB.UpdateTape(ctx, dbUpdateParams)
	if err != nil {
		return nil, err
	}

	tape := &model.Tape{
		ID:        dbTape.ID,
		PublicID:  dbTape.PublicID,
		CreatedAt: dbTape.CreatedAt,
		UpdatedAt: dbTape.UpdatedAt,
		Title:     dbTape.Title,
		Director:  dbTape.Director,
		Genre:     dbTape.Genre,
		Quantity:  dbTape.Quantity,
	}

	return tape, nil
}

func (r *tapeRepository) Delete(ctx context.Context, id int32) error {
	err := r.DB.DeleteTape(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *tapeRepository) DeleteAll(ctx context.Context) error {
	err := r.DB.DeleteAllTapes(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Helpers

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func toNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}
