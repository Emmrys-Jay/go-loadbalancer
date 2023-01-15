package domain

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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
	URL      *url.URL
	Proxy    *httputil.ReverseProxy
	Metadata map[string]string
}

func (s *Server) Forward(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

// GetMetaOrDefault returns the value associated with the given key
// in the metadata or returns a default.
func (s *Server) GetMetaOrDefault(key, def string) string {
	v, ok := s.Metadata[key]
	if !ok {
		return def
	}
	return v
}

// GetMetaOrDefaultInt returns the int value associated with the given key
// in the metadata or returns a default.
func (s *Server) GetMetaOrDefaultInt(key string, def int) int {
	v := s.GetMetaOrDefault(key, fmt.Sprintf("%d", def))
	a, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return a
}
