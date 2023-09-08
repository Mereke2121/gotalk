package server

import "net/http"

type WSServer struct {
	server *http.Server
}

func (s *WSServer) Run(port string, handler http.Handler) error {
	s.server = &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return s.server.ListenAndServe()
}
