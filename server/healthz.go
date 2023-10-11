package server

import (
	"net/http"

	"github.com/priyankasharma10/ReNew/utils"
	"github.com/sirupsen/logrus"
)

type healthResponse struct {
	Available bool `json:"up"`
}

func (srv *Server) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	logrus.WithField("test-key", "testing").WithField("test-key-2", "testing-2").Info("testing health route")
	utils.EncodeJSON200Body(w, healthResponse{
		Available: true})
}
