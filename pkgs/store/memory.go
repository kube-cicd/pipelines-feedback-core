package store

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Memory struct {
	mem map[string]string
}

func (m *Memory) Set(key, value string) error {
	logrus.Debugf("set(%s, %s)", key, value)
	m.mem[key] = value
	return nil
}

func (m *Memory) Get(key string) (string, error) {
	if val, ok := m.mem[key]; ok {
		return val, nil
	}
	return "", errors.New(ErrNotFound)
}

func NewMemory() *Memory {
	return &Memory{mem: make(map[string]string)}
}
