package backend

import "net/http"

func (cfg *ApiConfig) Reset(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	cfg.FileserverHits.And(0)
	resp.Write([]byte("Hits reset to 0\n"))
}
