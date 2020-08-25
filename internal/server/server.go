package server

import (
	"github.com/gorilla/mux"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
	store  interfaces.IStore
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func NewServer(store interfaces.IStore) *Server {
	server := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	server.ConfigureRouter()

	return server
}
