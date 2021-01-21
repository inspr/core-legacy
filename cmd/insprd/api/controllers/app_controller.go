package controller

import (
	ierrors "gitlab.inspr.dev/inspr/core/pkg/error"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// GetAllApps return all channels in the cluster's memory
func (s *Server) GetAllApps() {
}

// GetApp todo doc
func (s *Server) GetApp(query string) (*meta.App, error) {
	app, err := s.MemoryManager.Apps().GetApp(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't get a app in this path",
		).InternalServer().Build()
		return nil, serverErr
	}
	return app, nil
}

// CreateApp todo doc
func (s *Server) CreateApp(app *meta.App, ctx string) error {
	err := s.MemoryManager.Apps().CreateApp(app, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't create an app in this context and with these values",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// DeleteApp todo doc
func (s *Server) DeleteApp(query string) error {
	err := s.MemoryManager.Apps().DeleteApp(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't delete the App in this path",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// UpdateApp todo doc
func (s *Server) UpdateApp(app *meta.App, ctx string) error {
	err := s.MemoryManager.Apps().UpdateApp(app, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't modify the App in this context",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}
