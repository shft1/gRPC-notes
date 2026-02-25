package cors

import (
	"net/http"

	"github.com/shft1/grpc-notes/observability/logger"
)

func NewCORSMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			logger.Debug("cors headers have been successfully sets")

			next.ServeHTTP(w, r)
		})
	}
}
