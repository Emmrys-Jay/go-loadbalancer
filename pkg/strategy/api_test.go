package strategy_test

import (
	"testing"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/strategy"
)

func TestLoadStrategy(t *testing.T) {
	t.Run("Test RoundRobin", func(t *testing.T) {
		str := strategy.LoadStrategy("RoundRobin")

		_, ok := str.(*strategy.RoundRobin)
		if !ok {
			t.Errorf("Expected result strategy type to be 'RoundRobin'")
		}
	})

	t.Run("Test WeightedRoundRobin", func(t *testing.T) {
		str := strategy.LoadStrategy("WeightedRoundRobin")

		_, ok := str.(*strategy.WeightedRoundRobin)
		if !ok {
			t.Errorf("Expected result strategy type to be 'WeightedRoundRobin'")
		}
	})

	t.Run("Test Unknown", func(t *testing.T) {
		str := strategy.LoadStrategy("Unknown Type")

		_, ok := str.(*strategy.RoundRobin)
		if !ok {
			t.Errorf("Expected result strategy type to be 'RoundRobin'")
		}
	})
}
