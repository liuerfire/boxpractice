package httphandlers

import (
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

type loggingHandler struct {
	logger  logr.Logger
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	writer := &loggingWriter{writer: w}

	h.handler.ServeHTTP(writer, r)

	duration := time.Since(now)
	statusCode := writer.statusCode
	size := writer.size

	rawQuery := r.URL.RawQuery
	if rawQuery == "" {
		rawQuery = "-"
	}

	h.logger.Info("received HTTP request", "remoteAddr", r.RemoteAddr, "method", r.Method, "path", r.URL.Path, "query", rawQuery, "statusCode", statusCode, "size", size, "duration", duration.Nanoseconds(), "durationHuman", duration.String())
}

func LoggingHandler(logger logr.Logger) HandlerFunc {
	return func(h http.Handler) http.Handler {
		return loggingHandler{logger, h}
	}
}

type loggingWriter struct {
	writer     http.ResponseWriter
	statusCode int
	size       uint64
	committed  bool
}

func (w *loggingWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	if w.committed {
		return
	}
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
	w.committed = true
}

func (w *loggingWriter) Write(b []byte) (int, error) {
	if !w.committed {
		if w.statusCode == 0 {
			w.statusCode = http.StatusOK
		}
		w.WriteHeader(w.statusCode)
	}
	n, err := w.writer.Write(b)
	w.size += uint64(n)
	return n, err
}

func (w *loggingWriter) Flush() {
	w.writer.(http.Flusher).Flush()
}
