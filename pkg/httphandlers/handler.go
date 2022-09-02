package httphandlers

import (
	"net/http"
)

type HandlerFunc func(http.Handler) http.Handler

// Register registers a http.Handler
func Register(h http.Handler, hFns ...HandlerFunc) http.Handler {
	for _, hFn := range hFns {
		h = hFn(h)
	}
	return h
}
