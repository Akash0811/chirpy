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
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int    `json:expires_in_seconds`
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

	duration, err := getTimeDurationJWT(params.ExpiresIn)
	if err != nil {
		fmt.Printf("Failed to parse time due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, duration)
	if err != nil {
		fmt.Printf("Failed create JWT due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	type outgoingPayloadUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	payload := outgoingPayloadUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(
		resp,
		200,
		payload,
	)
}

func getTimeDurationJWT(expiresIn ...int) (time.Duration, error) {
	modifiedExpiresIn := defaultSeconds
	if len(expiresIn) > 1 {
		return time.Duration(0), fmt.Errorf("Expects at most one argument")
	} else if len(expiresIn) == 1 {
		if expiresIn[0] < defaultSeconds && expiresIn[0] > 0 {
			modifiedExpiresIn = expiresIn[0]
		}
	}
	parsedTime, err := time.ParseDuration(fmt.Sprintf("%vs", modifiedExpiresIn))
	if err != nil {
		return time.Duration(0), err
	}

	return parsedTime, nil
}
