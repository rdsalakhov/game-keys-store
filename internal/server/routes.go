package server

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/login", server.handleLogin()).Methods("GET")
	server.router.HandleFunc("/register", server.handleRegister()).Methods("POST")
	server.router.HandleFunc("/refresh", server.handleRefresh()).Methods("POST")

	server.router.HandleFunc("/game", server.authenticateSeller(server.handlePostGame())).Methods("POST")
	server.router.HandleFunc("/game/{id:[0-9]+}", server.handleFindGameByID()).Methods("GET")
	server.router.HandleFunc("/game", server.handleFindAllGames()).Methods("GET")
	server.router.HandleFunc("/game/{id:[0-9]+}", server.handleDeleteGameByID()).Methods("DELETE")
	server.router.HandleFunc("/game/{id:[0-9]+}/buy", server.handleBuyGame()).Methods("POST")
	server.router.HandleFunc("/game/{id:[0-9]+}/keys", server.authenticateSeller(server.handleGetGameKeys())).Methods("GET")

	server.router.HandleFunc("/key", server.authenticateSeller(server.handlePostKeys())).Methods("POST")

	server.router.HandleFunc("/purchase/{id:[0-9]+}", server.handlePostPurchase()).Methods("POST")
}
