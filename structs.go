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
	PolkaKey       string
}

type resErr struct {
	Error string `json:"error"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type email struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userReturnEmail struct {
	ID            uuid.UUID `json:"id"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	RefreshToken  string    `json:"refresh_token"`
	Is_Chirpy_Red bool      `json:"is_chirpy_red"`
}

type polkaRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"Data"`
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
