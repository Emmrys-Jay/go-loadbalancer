package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/config"
	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
	"github.com/Emmrys-Jay/go-loadbalancer/pkg/health"
	"github.com/Emmrys-Jay/go-loadbalancer/pkg/strategy"
)

var (
	port       = flag.Int("port", 8080, "port to start the load balancer")
	configPath = flag.String("config", "./example/config.yaml", "path to your config file")
)

type LoadBalancer struct {
	// Config is the configuration loaded from a config yaml file
	// TODO: This could be improved to fetch the configuration from
	// a more abstract concept (like ConfigSource) that can either be
	// a file or something else, and it should support hot reloading.
	Config *config.Config

	// the serverlist maps matcher to replicas
	ServerList map[string]*config.ServerList
}

func NewLoadBalancer(conf *config.Config) *LoadBalancer {
	// TODO: prevent multiple or invalid matchers before creating the server.
	serverList := make(map[string]*config.ServerList)

	for _, service := range conf.Services {
		servers := make([]*domain.Server, 0, len(service.Replicas))
		// Make all replicas into Servers
		for _, replica := range service.Replicas {
			url, err := url.Parse(replica.URL)
			if err != nil {
				log.Fatalln(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(url)
			servers = append(servers, &domain.Server{
				URL:      url,
				Proxy:    proxy,
				Metadata: replica.Metadata,
			})
		}
		checker, err := health.NewHealthChecker(nil, servers)
		if err != nil {
			log.Fatalln(err)
		}
		serverList[service.Matcher] = &config.ServerList{
			Servers:  servers,
			Name:     service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
			Hc:       checker,
		}
	}
	// starts all the health checkers for all matchers
	for _, v := range serverList {
		go v.Hc.Start()
	}
	return &LoadBalancer{
		Config:     conf,
		ServerList: serverList,
	}
}

// findServiceList looks for the first server that matches the reqPath (i.e matcher)
// will return an error if no matcher have been found
// TODO: Does it make sense to allow default responders?
func (l *LoadBalancer) findServiceList(reqPath string) (*config.ServerList, error) {
	log.Printf("Trying to find matcher for request '%s'", reqPath)
	for matcher, sl := range l.ServerList {
		if strings.HasPrefix(reqPath, matcher) {
			log.Infof("Found service '%s' matching the request", sl.Name)
			return sl, nil
		}
	}

	return nil, fmt.Errorf("could not find a matcher for url '%s'", reqPath)
}

func (l *LoadBalancer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// TODO: we need to support per service forwarding, i.e this method
	// should read request path, say host:port/service/rest/of/url , this should be
	// forwarded against service named "service" and url will be "host(i):port(i)/rest/of/url"
	log.Infof("Received new request: url=%s", r.Host)

	sl, err := l.findServiceList(r.URL.Path)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := sl.Strategy.Next(sl.Servers)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("Forwarding to server number %s", server.URL.Host)
	log.Info("\n")
	// Forwarding the request to the proxy
	server.Forward(rw, r)
}

func main() {
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("while loading config, got error: '%v'\n", err)
	}

	lb := NewLoadBalancer(conf)

	srv := http.Server{
		Addr:    ":" + fmt.Sprint(*port),
		Handler: lb,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
