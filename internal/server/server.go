package server

import (
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/service"
)

type Server struct {
	connectorService *service.Connector
}

func NewServer(connectorService *service.Connector) *Server {
	return &Server{
		connectorService: connectorService,
	}
}

func (s *Server) Start() error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}
