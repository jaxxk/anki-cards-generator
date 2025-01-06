package create

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureAnkiConnect(t *testing.T) {
	ok, err := EnsureAnkiConnect()
	assert.NoError(t, err, "expected no error, but got one")
	if !ok {
		t.Fatal("expected true but got false")
	}
}
