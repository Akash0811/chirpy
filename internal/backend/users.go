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

func (cfg *ApiConfig) AddUser(resp http.ResponseWriter, req *http.Request) {
	type userReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := userReq{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Failed to decode input due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	type outgoingPayloadUser struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		fmt.Printf("Failed to hash password due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}
	user, err := cfg.Queries.CreateUser(
		req.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		fmt.Printf("Failed to insert user in database %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	payload := outgoingPayloadUser{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(
		resp,
		201,
		payload,
	)
}

func (cfg *ApiConfig) UpdateUser(resp http.ResponseWriter, req *http.Request) {
	type userReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := userReq{}
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

	user, err := cfg.Queries.GetUser(req.Context(), userId)
	if err != nil {
		fmt.Printf("User not found due to %v\n", err)
		respondWithError(resp, 401, "Token invalid/expired")
		return
	}

	type outgoingPayloadUser struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red "`
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		fmt.Printf("Failed to hash password due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}
	user, err = cfg.Queries.UpdateUserDetails(
		req.Context(),
		database.UpdateUserDetailsParams{
			ID:             user.ID,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		fmt.Printf("Failed to update user in database %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	payload := outgoingPayloadUser{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(
		resp,
		200,
		payload,
	)
}

func (cfg *ApiConfig) UpdateUserMembership(resp http.ResponseWriter, req *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}
	type webhookReq struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := webhookReq{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Failed to decode input due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(resp, 204, struct{}{})
		return
	}

	apiKey, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithError(resp, 401, "Api key invalid")
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(resp, 401, "Api key invalid")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(resp, 500, serverErrorString)
		return
	}
	err = cfg.Queries.UpdateUserChrirpyRed(
		req.Context(),
		database.UpdateUserChrirpyRedParams{
			ID:          userID,
			IsChirpyRed: true,
		},
	)
	if err != nil {
		fmt.Printf("Failed to update user in database %v\n", err)
		respondWithError(resp, 404, "User not found")
		return
	}

	respondWithJSON(
		resp,
		204,
		struct{}{},
	)
}
