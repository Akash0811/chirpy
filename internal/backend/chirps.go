package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Akash0811/chirpy/internal/auth"
	"github.com/Akash0811/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) AddChirp(resp http.ResponseWriter, req *http.Request) {
	type validateChirpIncomingPayload struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := validateChirpIncomingPayload{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Failed to decode input due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, 401, "Token invalid/expired")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		fmt.Printf("Failed to validate token due to %v\n", err)
		respondWithError(resp, 401, "Token invalid/expired")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	user, err := cfg.Queries.GetUser(req.Context(), userId)
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

func (cfg *ApiConfig) GetAllChirps(resp http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.Queries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(resp, 404, "User not found")
		return
	}

	type payloadChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	var payload []payloadChirp

	for _, dbChirp := range dbChirps {
		payload = append(payload, payloadChirp(dbChirp))
	}
	respondWithJSON(resp, 200, payload)
}

func (cfg *ApiConfig) GetChirp(resp http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		fmt.Printf("Could not parse uuid %v\n", req.PathValue("chirpID"))
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	dbChirp, err := cfg.Queries.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(resp, 404, "User not found")
		return
	}

	type payloadChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(resp, 200, payloadChirp(dbChirp))
}
