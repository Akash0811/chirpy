package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Akash0811/chirpy/internal/auth"
	"github.com/Akash0811/chirpy/internal/database"
)

func (cfg *ApiConfig) RefreshToken(resp http.ResponseWriter, req *http.Request) {
	refreshtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, 401, "Token invalid/expired")
		return
	}
	refreshToken, err := cfg.Queries.GetRefreshToken(req.Context(), refreshtoken)
	if err != nil {
		fmt.Printf("No tokens found due to %v\n", err)
		respondWithError(resp, 401, "Refresh Token not found or revoked or expired")
		return
	}
	currentTime := time.Now()
	flagRevoked := false
	if refreshToken.RevokedAt.Valid {
		if refreshToken.RevokedAt.Time.Before(currentTime) {
			flagRevoked = true
		}
	}
	if flagRevoked || refreshToken.ExpiresAt.Before(currentTime) {
		fmt.Printf("Tokens already revoked or expired %v\n", err)
		respondWithError(resp, 401, "Refresh Token not found or revoked or expired")
		return
	}

	err = cfg.Queries.UpdateRefreshToken(
		req.Context(),
		database.UpdateRefreshTokenParams{
			Token:     refreshToken.Token,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		fmt.Printf("Failed to update token details due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	duration, err := getTimeDurationJWT()
	if err != nil {
		fmt.Printf("Failed to parse time due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	jwttoken, err := auth.MakeJWT(refreshToken.UserID, cfg.JWTSecret, duration)
	if err != nil {
		fmt.Printf("Failed create JWT due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	type outgoingPayloadUser struct {
		JWTToken string `json:"token"`
	}
	payload := outgoingPayloadUser{
		JWTToken: jwttoken,
	}
	respondWithJSON(
		resp,
		200,
		payload,
	)
}

func (cfg *ApiConfig) RevokeToken(resp http.ResponseWriter, req *http.Request) {
	refreshtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resp, 401, "Token invalid/expired")
		return
	}
	refreshToken, err := cfg.Queries.GetRefreshToken(req.Context(), refreshtoken)
	if err != nil {
		fmt.Printf("No tokens found due to %v\n", err)
		respondWithError(resp, 401, "Refresh Token not found or revoked or expired")
		return
	}
	currentTime := time.Now()
	flagRevoked := false
	if refreshToken.RevokedAt.Valid {
		if refreshToken.RevokedAt.Time.Before(currentTime) {
			flagRevoked = true
		}
	}
	if flagRevoked || refreshToken.ExpiresAt.Before(currentTime) {
		fmt.Printf("Tokens already revoked or expired %v\n", err)
		respondWithError(resp, 401, "Refresh Token not found or revoked or expired")
		return
	}

	err = cfg.Queries.RevokeRefreshToken(
		req.Context(),
		database.RevokeRefreshTokenParams{
			Token:     refreshToken.Token,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		fmt.Printf("Failed to update token details due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	respondWithJSON(
		resp,
		204,
		struct{}{},
	)
}
