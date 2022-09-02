package httphandlers

import (
	"net/http"
)

type corsHandler struct {
	handler http.Handler
}

// Added this to allow us to send requests using Swagger UI.
// Only for development.
func (c corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		c.handler.ServeHTTP(w, r)
	}
}

func CorsConfigHandler() HandlerFunc {
	return func(handler http.Handler) http.Handler {
		return corsHandler{
			handler,
		}
	}
}
