package store

import "github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"

const ErrNotFound = "No such key"

type Store interface {
	contract.Pluggable
	Set(key string, value string, ttl int) error
	Get(key string) (string, error)
	Initialize() error
}
