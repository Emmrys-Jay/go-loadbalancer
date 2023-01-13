package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/config"
)

var (
	port       = flag.Int("port", 8080, "port to start the load balancer")
	configPath = flag.String("config", "./example/config.yaml", "path to your config file")
)

type LoadBalancer struct {
	Config     *config.Config
	ServerList *config.ServerList
}

func NewLoadBalancer(conf *config.Config) *LoadBalancer {
	servers := make([]*config.Server, 0, len(conf.Services))
	for _, service := range conf.Services {
		// TODO: Don't ignore names
		for _, replica := range service.Replicas {
			url, err := url.Parse(replica)
			if err != nil {
				log.Fatalln(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(url)
			servers = append(servers, &config.Server{
				URL:   url,
				Proxy: proxy,
			})
		}
	}

	return &LoadBalancer{
		Config: conf,
		ServerList: &config.ServerList{
			Servers: servers,
			Current: uint32(0),
		},
	}
}

func (l *LoadBalancer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// TODO: we need to support per service forwarding, i.e this method
	// should read request path, say host:port/service/rest/of/url , this should be
	// forwarded against service named "service" and url will be "host(i):port(i)/rest/of/url"
	log.Printf("Received new request: url=%s", r.Host)

	next := l.ServerList.Next()
	log.Printf("Forwarding to server number '%d'", next)
	// Forwarding the request to the proxy
	l.ServerList.Servers[next].Forward(rw, r)
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
