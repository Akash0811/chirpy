package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) AddUser(resp http.ResponseWriter, req *http.Request) {
	type userReq struct {
		Email string `json:"email"`
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
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	user, err := cfg.Queries.CreateUser(req.Context(), params.Email)
	if err != nil {
		fmt.Printf("Failed to insert user in database %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	payload := outgoingPayloadUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(
		resp,
		201,
		payload,
	)
}
