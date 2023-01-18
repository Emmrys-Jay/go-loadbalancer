package strategy

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
)

const (
	S_RoundRobin         = "RoundRobin"
	S_WeightedRoundRobin = "WeightedRoundRobin"
	S_Unknown            = "Unknown"
)

var strategies = map[string]BalancingStrategy{
	S_RoundRobin:         &RoundRobin{current: 0, mu: sync.Mutex{}},
	S_WeightedRoundRobin: &WeightedRoundRobin{mu: sync.Mutex{}},
}

// BalancingStrategy is the abstraction that allow for using different
// load balancing strategies.
type BalancingStrategy interface {
	Next(servers []*domain.Server) (*domain.Server, error)
}

// RoundRobin implements BalancingStrategy
type RoundRobin struct {
	mu sync.Mutex
	// the current server to forward the request to.
	// the next server should be (current + 1) % len(servers)
	current int
}

func (rr *RoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	seen := 0
	for seen < len(servers) {
		rr.current = (rr.current + 1) % len(servers)

		if servers[rr.current].IsAlive() {
			return servers[rr.current], nil
		}
		seen++
	}

	log.Warnf("All servers are down")
	return nil, errors.New("all servers are down")
}

// WeightedRoundRobin is a strategy that is similar to the RoundRobin, but
// it takes server compute power into consideration. The compute power of a server is
// given as an integer representing the fraction of requests that one server can handle
// over another.
//
// RoundRobin strategy is equivalent to a WeightedRoundRobin with all weights = 1
type WeightedRoundRobin struct {
	// Any changes to the below fields must be done while holding the 'mu' lock.
	mu sync.Mutex
	// Note: This is making the assumption that the server lists coming through the
	// Next function won't change between successive calls.
	// Changing the server list would cause this strategy to break, panic, or
	// not route properly.
	//
	// count will keep track of the number of request server 'i' has processed.
	count []int
	// current is the index of the last server that executed a request.
	current int
}

func (wrr *WeightedRoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()
	seen := 0
	if wrr.count == nil {
		// First time using the strategy
		wrr.count = make([]int, len(servers))
		wrr.current = 0
	}

	capacity := servers[wrr.current].GetMetaOrDefaultInt("weight", 1)

	for seen < len(servers) {
		if wrr.count[wrr.current] < capacity && servers[wrr.current].IsAlive() {
			// Current server can still accept some requests
			wrr.count[wrr.current]++
			log.Printf("Strategy picked server %s", servers[wrr.current].URL.Host)
			return servers[wrr.current], nil
		}
		seen++

		// Server has gotten to its limit
		// Reset the current one and move to the next
		wrr.count[wrr.current] = 0
		wrr.current = (wrr.current + 1) % len(servers)
	}

	log.Warn("All servers are down")
	return nil, errors.New("all servers are down")
}

// LoadStrategy tries to resolve the balancing strategy based on its name.
// It defaults to RoundRobin if an unknown name is passed in.
func LoadStrategy(name string) BalancingStrategy {
	if st, ok := strategies[name]; ok {
		log.Infof("Picked strategy '%s'", name)
		return st
	}
	log.Warnf("Look up on strategy '%s' unsuccessful, defaulting to 'RoundRobin'", name)
	return strategies[S_RoundRobin]
}
