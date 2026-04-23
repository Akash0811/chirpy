package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Akash0811/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) AddChirp(resp http.ResponseWriter, req *http.Request) {
	type validateChirpIncomingPayload struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := validateChirpIncomingPayload{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Failed to decode input due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(resp, 400, inputVaildationErrorString)
	}

	user, err := cfg.Queries.GetUser(req.Context(), params.UserID)
	if err != nil {
		respondWithError(resp, 404, "User not found")
		return
	}

	cleanedBody, err := replaceBadWords(params.Body)
	if err != nil {
		fmt.Printf("Failed to clean body due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	chirp, err := cfg.Queries.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: user.ID,
		},
	)
	type payload struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(resp, 201, payload{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
