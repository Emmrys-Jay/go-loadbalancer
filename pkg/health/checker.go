package health

import (
	"errors"
	"net"
	"time"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
	log "github.com/sirupsen/logrus"
)

type HealthChecker struct {
	servers []*domain.Server
	// TODO: configure the period based on the config file
	period int
}

// NewHealthChecker will create a new health checker
func NewHealthChecker(_conf *domain.Config, servers []*domain.Server) (*HealthChecker, error) {
	if len(servers) == 0 {
		return nil, errors.New("cannot start checkup for empty server list")
	}
	return &HealthChecker{
		servers: servers,
		period:  1,
	}, nil
}

// Start attempts to indefinitely check the health of each server
// The caller is responsible for creating the goroutine where this should run
func (hc *HealthChecker) Start(name string) {
	log.Infof("Starting the health checks for service '%s'...", name)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		for _, server := range hc.servers {
			go checkHealth(server)
		}
	}
}

// checkHealth should checks if a server is live and changes it's alive field
func checkHealth(server *domain.Server) {
	// we will consider a server to be healthy if we can open a tcp connection
	// to the host:port within reasonable time frame
	_, err := net.DialTimeout("tcp", server.URL.Host, time.Second*30)
	if err != nil {
		log.Errorf("Could not connect to server at %s", server.URL.Host)
		if old := server.SetLiveness(false); old {
			log.Warnf("Transitioning server '%s' from 'Healthy' to 'Unhealthy'", server.URL.Host)
		}
		return
	}

	if old := server.SetLiveness(true); !old {
		log.Infof("Transitioning server '%s' from 'Unhealthy' to 'Healthy'", server.URL.Host)
	}
}
