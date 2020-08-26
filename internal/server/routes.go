package server

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/login", server.handleLogin()).Methods("GET")
	server.router.HandleFunc("/register", server.handleRegister()).Methods("POST")
	server.router.HandleFunc("/refresh", server.handleRefresh()).Methods("POST")
}
