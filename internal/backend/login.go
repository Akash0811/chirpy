package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Akash0811/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) LoginUser(resp http.ResponseWriter, req *http.Request) {
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

	user, err := cfg.Queries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		fmt.Printf("No users found due to %v\n", err)
		respondWithError(resp, 404, "User not found")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Something wrong when verifying password due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}
	if !match {
		fmt.Printf("Password Verification failed due to %v\n", err)
		respondWithError(resp, 401, "Password verification failed")
		return
	}

	type outgoingPayloadUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	payload := outgoingPayloadUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(
		resp,
		200,
		payload,
	)
}
