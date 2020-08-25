package server

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	ctxKeyRequestID ctxKey = iota
)

type ctxKey int

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}
