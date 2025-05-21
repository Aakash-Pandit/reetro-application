package server

import (
	"github.com/Aakash-Pandit/reetro-golang/routes"
)

type ServerApplication struct {
	APIRoute *routes.Router
}

func NewServer(route *routes.Router) *ServerApplication {
	return &ServerApplication{
		APIRoute: route,
	}
}

func (s *ServerApplication) Start() {
	s.APIRoute.Run()
}
