package server

import "net/http"

func (s *server) bindRoutes() {
	routes := map[string]http.HandlerFunc{
		"/upload": s.handleUpload(),
		"/":       s.handleIndex(),
	}

	for path, handler := range routes {
		s.router.Handle(path, s.LogRequest(handler))
	}
}
