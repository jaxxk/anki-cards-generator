package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	client, err := newClient(zap.NewExample().Sugar())
	assert.NoError(t, err)
	if client == nil {
		t.Fatal("expected client to never be nil")
	}
}
