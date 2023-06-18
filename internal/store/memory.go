package store

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	"time"
)

type MemEntry struct {
	val     string
	expires time.Time
}

type Memory struct {
	mem map[string]MemEntry
}

func (m *Memory) Set(key, value string, ttl int) error {
	if ttl == 0 {
		ttl = 86400 * 365 * 10 // 10 years should be enough
	}
	expires := time.Now()
	expires.Add(time.Second * time.Duration(ttl))

	m.mem[key] = MemEntry{
		val:     value,
		expires: expires,
	}
	return nil
}

func (m *Memory) Get(key string) (string, error) {
	if entry, ok := m.mem[key]; ok {
		if entry.expires.After(time.Now()) {
			return "", errors.New(store.ErrNotFound)
		}
		return entry.val, nil
	}
	return "", errors.New(store.ErrNotFound)
}

func NewMemory() *Memory {
	return &Memory{mem: make(map[string]MemEntry)}
}
