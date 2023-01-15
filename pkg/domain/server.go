package domain

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Replica struct {
	URL      string            `yaml:"url"`
	Metadata map[string]string `yaml:"metadata"`
}

type Service struct {
	Name string `yaml:"name"`

	// A prefix matcher to select service based on the path part of the url
	Matcher string `yaml:"matcher"`

	// Strategy is the load balancing strategy for the current service
	Strategy string    `yaml:"strategy"`
	Replicas []Replica `yaml:"replicas"`
}

// Server is an instance of a running server
type Server struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
}

func (s *Server) Forward(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}
