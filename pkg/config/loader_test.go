package config_test

import (
	"testing"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/config"
)

func TestLoadConfig(t *testing.T) {
	conf, err := config.LoadConfig("./testdata/test.yaml")
	if err != nil {
		t.Errorf("Expected error to be nil, got '%v'", err)
	}

	if conf.Strategy != "RoundRobin" {
		t.Errorf("Expected strategy 'RoundRobin', got '%s'", conf.Strategy)
	}

	if len(conf.Services) != 1 {
		t.Errorf("Expected servcies count to be '1', got '%d'", len(conf.Services))
	}

	if conf.Services[0].Name != "Test service 1" {
		t.Errorf("Expected service name to be  'Test service 1', got '%s'", conf.Services[0].Name)
	}

	if conf.Services[0].Strategy != "RoundRobin" {
		t.Errorf("Expected strategy for service 'RoundRobin', got '%s'", conf.Strategy)
	}

	if conf.Services[0].Matcher != "/api/v1" {
		t.Errorf("Expected matcher '/api/v1', got '%s'", conf.Services[0].Matcher)
	}

	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("Expected service replicas count to be '2', got '%d'", len(conf.Services[0].Replicas))
	}

	if conf.Services[0].Replicas[0] != "http://localhost:8081" {
		t.Errorf("Expected service replica 1 to be  'http://localhost:8081', got '%s'", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1] != "http://localhost:8082" {
		t.Errorf("Expected service name to be  'http://localhost:8082', got '%s'", conf.Services[0].Replicas[1])
	}
}
