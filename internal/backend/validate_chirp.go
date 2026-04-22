package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidateChirp(resp http.ResponseWriter, req *http.Request) {
	type validateChirpIncomingPayload struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := validateChirpIncomingPayload{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(resp, 500, fmt.Sprintf("Failed to decode incoming data due to %v", err))
		return
	}

	type outgoingPayload struct {
		Valid bool `json:"valid"`
	}

	if len(params.Body) <= 140 {
		respondWithJSON(resp, 200, outgoingPayload{Valid: true})
	} else {
		respondWithError(resp, 400, "Chirp is too long")
	}

}
