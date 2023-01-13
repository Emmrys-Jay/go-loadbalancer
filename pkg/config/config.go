package config

import (
	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
	"github.com/Emmrys-Jay/go-loadbalancer/pkg/strategy"
)

// Config is a representation of the configuration
// given to the LB from a config source
type Config struct {
	Services []domain.Service `yaml:"services"`

	// Strategy is the name of strategy to be used in load balancing between instances
	Strategy string `yaml:"strategy"`
}

type ServerList struct {
	// Servers are the replicas
	Servers []*domain.Server

	// Name is the name of the service
	Name string

	// Strategy defines how the server list is load balanced
	// It can never be nil, instead it defaults to RoundRobin
	Strategy strategy.BalancingStrategy
}
