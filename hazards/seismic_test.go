package hazards_test

import (
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

func TestMagnitude(t *testing.T) {
	d := hazards.SeismicEvent{}
	d.SetDepth(2.5)
	if d.Depth() != 2.5 {
		t.Errorf("Expected %f, got %f", 2.5, d.Depth())
	}
}
