package contract

type Pluggable interface {
	CanHandle(adapterName string) bool
}
