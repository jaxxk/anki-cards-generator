package transform

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	client := newClient()
	if client == nil {
		t.Fatal("expected client to never be nil")
	}
}
