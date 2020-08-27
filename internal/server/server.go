package server

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
	store  interfaces.IStore
	redis  *redis.Client
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func NewServer(store interfaces.IStore, redis *redis.Client) *Server {
	server := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		redis:  redis,
	}

	server.ConfigureRouter()

	return server
}

func (server *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	server.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (server *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
