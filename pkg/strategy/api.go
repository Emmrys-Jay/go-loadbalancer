package strategy

import (
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
)

const (
	S_RoundRobin         = "RoundRobin"
	S_WeightedRoundRobin = "WeightedRoundRobin"
	S_Unknown            = "Unknown"
)

var strategies = map[string]BalancingStrategy{
	S_RoundRobin:         &RoundRobin{current: uint32(0)},
	S_WeightedRoundRobin: &WeightedRoundRobin{mu: sync.Mutex{}},
}

// BalancingStrategy is the abstraction that allow for using different
// load balancing strategies.
type BalancingStrategy interface {
	Next(servers []*domain.Server) *domain.Server
}

// RoundRobin implements BalancingStrategy
type RoundRobin struct {
	// the current server to forward the request to.
	// the next server should be (current + 1) % len(servers)
	current uint32
}

func (rr *RoundRobin) Next(servers []*domain.Server) *domain.Server {
	nxt := atomic.AddUint32(&rr.current, uint32(1))
	lenS := uint32(len(servers))
	picked := servers[nxt%lenS]
	log.Printf("Strategy picked server %s", picked.URL.Host)
	return picked
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
	// Changing the server list would cause this strategy tp break, panic, or
	// not route properly.
	//
	// count will keep track of the number of request server 'i' has processed.
	count []int
	// current is the index of the last server that executed a request.
	current int
}

func (wrr *WeightedRoundRobin) Next(servers []*domain.Server) *domain.Server {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()
	if wrr.count == nil {
		// First time using the strategy
		wrr.count = make([]int, len(servers))
		wrr.current = 0
	}

	capacity := servers[wrr.current].GetMetaOrDefaultInt("weight", 1)
	if wrr.count[wrr.current] < capacity {
		// Current server can still accept some requests
		wrr.count[wrr.current]++
		log.Printf("Strategy picked server %s", servers[wrr.current].URL.Host)
		return servers[wrr.current]
	}

	// Server has gotten to its limit
	// Reset the current one and move to the next
	wrr.count[wrr.current] = 0
	wrr.current = (wrr.current + 1) % len(servers)
	log.Printf("Strategy picked server %s", servers[wrr.current].URL.Host)
	return servers[wrr.current]
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
