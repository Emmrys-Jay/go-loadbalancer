package domain_test

import (
	"testing"

	"github.com/Emmrys-Jay/go-loadbalancer/pkg/domain"
)

type testData struct {
	name        string
	firstValue  bool
	secondValue bool
	expected    bool
}

var data = []testData{
	{"TestSetLiveness_False-True", false, true, false},
	{"TestSetLiveness_False-False", false, false, false},
	{"TestSetLiveness_True-False", true, false, true},
	{"TestSetLiveness_True-True", true, true, true},
}

func TestSetLiveness(t *testing.T) {
	s := domain.Server{}
	go func() {
		r := s.SetLiveness(true)
		if r {
			t.Errorf("Expected server liveness to be true")
		}
	}()

	for _, val := range data {
		t.Run(val.name, func(t *testing.T) {
			go func() {
				server := domain.Server{}
				_ = server.SetLiveness(val.firstValue)
				result := server.SetLiveness(val.secondValue)

				if result != val.expected {
					t.Errorf("Expected server liveness to be '%v'", val.expected)
				}
			}()
		})
	}

}
