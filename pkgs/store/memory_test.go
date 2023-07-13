package store_test

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemory_SetGet_WithTTL(t *testing.T) {
	m := store.NewMemory()
	m.Set("book", "conquest-of-bread", 1)

	retrieved, _ := m.Get("book")
	assert.Equal(t, "conquest-of-bread", retrieved)

	// WARNING: This is a time-based test
	time.Sleep(time.Second * 2)

	// after two seconds the entry should no longer exist, as it has TTL=1s
	retrieved, _ = m.Get("book")
	assert.Equal(t, "", retrieved)
}

func TestMemory_SetGet_WithoutTTL(t *testing.T) {
	m := store.NewMemory()
	m.Set("book", "conquest-of-bread", 0)

	retrieved, _ := m.Get("book")
	assert.Equal(t, "conquest-of-bread", retrieved)

	// WARNING: This is a time-based test
	time.Sleep(time.Second * 1)

	// TTL=0s, so the entry should still exist
	retrieved, _ = m.Get("book")
	assert.Equal(t, "conquest-of-bread", retrieved)
}
