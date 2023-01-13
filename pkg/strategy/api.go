package strategy

import (
	"log"
	"sync/atomic"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
)

const (
	S_RoundRobin         = "RoundRobin"
	S_WeightedRoundRobin = "WeightedRoundRobin"
	S_Unknown            = "Unknown"
)

var strategies = map[string]BalancingStrategy{
	S_RoundRobin: &RoundRobin{current: uint32(0)},
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

// LoadStrategy tries to resolve the balancing strategy based on its name.
// It defaults to RoundRobin if an unknown name is passed in.
func LoadStrategy(name string) BalancingStrategy {
	if st, ok := strategies[name]; ok {
		return st
	}

	return strategies[S_RoundRobin]
}
