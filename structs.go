package main

import (
	"sync/atomic"
	"time"

	"github.com/christianrm0821/Chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	PLATFORM       string
	Secret         string
}

type resErr struct {
	Error string `json:"error"`
}

type email struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

type userReturnEmail struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type chirpPostReq struct {
	Body string `json:"body"`
	// UserID uuid.UUID `json:"user_id"`
}

type validChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
