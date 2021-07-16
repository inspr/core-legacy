package controller

import (
	"go.uber.org/zap"
	handler "inspr.dev/inspr/cmd/uid_provider/api/handlers"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest"
)

// initRoutes defines which routes the UID Provider API will have
func (s *Server) initRoutes() {
	_, alevel := logs.Logger(zap.Fields(zap.String("section", "uidp-server")))
	logger.Debug("initializing UIDP server routes")

	h := handler.NewHandler(s.ctx, s.rdb)

	s.mux.HandleFunc("/newuser", h.CreateUserHandler())

	s.mux.HandleFunc("/deleteuser", h.DeleteUserHandler())

	s.mux.HandleFunc("/updatepwd", h.UpdatePasswordHandler())

	s.mux.HandleFunc("/login", h.LoginHandler())

	s.mux.HandleFunc("/refreshtoken", h.RefreshTokenHandler())

	s.mux.Handle("/healthz", rest.Healthz())

	s.mux.Handle("/log/level", alevel)
}
