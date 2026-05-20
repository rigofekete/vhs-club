package handler

import (
	"time"

	"github.com/google/uuid"
)

type RentalResponse struct {
	PublicID  uuid.UUID `json:"public_id"`
	TapeID    int32     `json:"tape_id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	TapeTitle string    `json:"tape_title"`
	Username  string    `json:"username"`
	RentedAt  time.Time `json:"rented_at"`
}
