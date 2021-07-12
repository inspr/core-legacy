package controller

import handler "inspr.dev/inspr/cmd/uid_provider/api/handlers"

// initRoutes defines which routes the UID Provider API will have
func (s *Server) initRoutes() {
	logger.Debug("initializing Insprd server routes")

	h := handler.NewHandler(s.ctx, s.rdb)

	s.mux.HandleFunc("/newuser", h.CreateUserHandler())

	s.mux.HandleFunc("/deleteuser", h.DeleteUserHandler())

	s.mux.HandleFunc("/updatepwd", h.UpdatePasswordHandler())

	s.mux.HandleFunc("/login", h.LoginHandler())

	s.mux.HandleFunc("/refreshtoken", h.RefreshTokenHandler())
}
