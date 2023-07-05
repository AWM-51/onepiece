package Callutil

import (
	"context"
)

type Caller interface {
	SetBeforeRequestAspect(f func(context.Context, string, string, interface{}) error) Caller
	SetAfterResponseAspect(f func(context.Context, []byte, int) error) Caller
	Call(req interface{}) (*Response, error)
}
