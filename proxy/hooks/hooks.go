package hooks

import (
	"net/http"
)

// PreRequest type
type PreRequest func(*http.Request) error

// PostRequest type
type PostRequest func() func(*http.Response) error

// Register interface
type Register interface {
	HookBefore()
	HookAfter()
}

// New func
func New() Register {
	return &register{}
}

type register struct {
}

func (r *register) HookBefore() {

}

func (r *register) HookAfter() {

}
