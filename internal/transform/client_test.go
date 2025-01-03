package transform

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	client := newClient()
	if client == nil {
		t.Fatal("expected client to never be nil")
	}
}
func TestNewCompletion(t *testing.T) {
	t.Parallel()
	client := NewChatCompletion(context.TODO(), "")
	if client == nil {
		t.Fatal("expected completion to never be nil")
	}
}
