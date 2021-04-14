package controller

import handler "github.com/inspr/inspr/cmd/uid_provider/api/handlers"

// initRoutes defines which routes the UID Provider API will have
func (s *Server) initRoutes() {
	h := handler.NewHandler(s.Rdb)

	s.Mux.HandleFunc("/newuser", h.CreateUserHandler)

	s.Mux.HandleFunc("/deleteuser", h.DeleteUserHandler)

	s.Mux.HandleFunc("/updatepwd", h.UpdatePasswordHandler)

	s.Mux.HandleFunc("/login", h.LoginHandler)

	s.Mux.HandleFunc("/refreshtoken", h.RefreshTokenHandler)
}
