package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type server struct {
	port   string
	router *http.ServeMux
	logger *logrus.Logger
}

func New() (*server, error) {
	portEnv := os.Getenv("APP_PORT")
	_, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, err
	}

	s := server{
		port:   portEnv,
		logger: logrus.New(),
		router: http.NewServeMux(),
	}

	s.bindRoutes()

	return &s, nil
}

func (s *server) Serve() error {
	s.logger.Info("Starting server on port " + s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
