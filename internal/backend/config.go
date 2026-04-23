package backend

import (
	"sync/atomic"

	"github.com/Akash0811/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Queries        *database.Queries
	Platform       string
}
