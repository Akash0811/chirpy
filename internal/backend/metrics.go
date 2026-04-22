package backend

import (
	"fmt"
	"net/http"
)

const metricsPageContent = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *ApiConfig) Metrics(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	numHits := cfg.FileserverHits.Add(0)
	resp.Write([]byte(fmt.Sprintf(metricsPageContent, numHits)))
}
