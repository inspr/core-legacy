package controller

import handler "gitlab.inspr.dev/inspr/core/cmd/uid_provider/api/handlers"

func (s *Server) initRoutes() {
	h := handler.NewHandler(s.Rdb)

	s.Mux.HandleFunc("/newuser", h.CreateUserHandler)

	s.Mux.HandleFunc("/deleteuser", h.CreateUserHandler)

	s.Mux.HandleFunc("/updatepwd", h.CreateUserHandler)

	s.Mux.HandleFunc("/login", h.CreateUserHandler)

	s.Mux.HandleFunc("/refreshtoken", h.CreateUserHandler)
}
