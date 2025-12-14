package input

import (
	"testing"
)

func TestNewEbitenInputAdapter(t *testing.T) {
	if NewEbitenInputAdapter() == nil {
		t.Fatalf("expected non-nil adapter")
	}
}
