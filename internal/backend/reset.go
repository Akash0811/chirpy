package backend

import (
	"net/http"
)

func (cfg *ApiConfig) Reset(resp http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(resp, 403, "Forbidden Action")
		return
	}

	err := cfg.Queries.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(resp, 500, serverErrorString)
		return
	}
	cfg.FileserverHits.And(0)

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte("Hits reset to 0\nDeleted All users\n"))
}
