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
		fmt.Printf("Failed to decode input due to %v\n", err)
		respondWithError(resp, 400, inputVaildationErrorString)
		return
	}

	type outgoingPayload struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cleanedBody, err := replaceBadWords(params.Body)
	if err != nil {
		fmt.Printf("Failed to clean body due to %v\n", err)
		respondWithError(resp, 500, serverErrorString)
		return
	}

	if len(params.Body) <= 140 {
		respondWithJSON(resp, 200, outgoingPayload{CleanedBody: cleanedBody})
	} else {
		respondWithError(resp, 400, "Chirp is too long")
	}

}
