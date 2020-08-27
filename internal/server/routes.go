package server

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/login", server.handleLogin()).Methods("GET")
	server.router.HandleFunc("/register", server.handleRegister()).Methods("POST")
	server.router.HandleFunc("/refresh", server.handleRefresh()).Methods("POST")

	//game := server.router.PathPrefix("/game").Subrouter()
	//game.Use(server.authenticateSeller)
	server.router.HandleFunc("/game", server.authenticateSeller(server.handlePostGame())).Methods("POST")
}
