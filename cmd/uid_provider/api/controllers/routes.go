package controller

import handler "gitlab.inspr.dev/inspr/core/cmd/uid_provider/api/handlers"

func (s *Server) initRoutes() {
	h := handler.NewHandler(s.Rdb)

	s.Mux.HandleFunc("/newuser", h.CreateUserHandler)

	s.Mux.HandleFunc("/deleteuser", h.DeleteUserHandler)

	s.Mux.HandleFunc("/updatepwd", h.UpdatePasswordHandler)

	s.Mux.HandleFunc("/login", h.LoginHandler)

	s.Mux.HandleFunc("/refreshtoken", h.RefreshTokenHandler)
}
