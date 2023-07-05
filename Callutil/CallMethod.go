package Callutil

import (
	"context"
	"errors"
	"time"
)

type Client struct {
	Protocol string
	Timeout  time.Duration
	callFunc func(interface{}) (*Response, error)
}

type CallMethod struct {
	Client
	Url                 string
	Method              string
	Headers             map[string]string
	Params              map[string]string
	BeforeRequestAspect func(context.Context, string, string, interface{}) error
	AfterResponseAspect func(context.Context, []byte, int) error
}

// NewCallMethod creates a new CallMethod instance.
func NewCallMethod() *CallMethod {
	return &CallMethod{}
}

func (c *Client) SetProtocol(protocol string) *Client {
	c.Protocol = protocol
	return c
}

func (c *CallMethod) SetParams(params map[string]string) *CallMethod {
	c.Params = params
	return c
}
func (c *CallMethod) SetBeforeRequestAspect(f func(context.Context, string, string, interface{}) error) Caller {
	c.BeforeRequestAspect = f
	return c
}

func (c *CallMethod) SetAfterResponseAspect(f func(context.Context, []byte, int) error) Caller {
	c.AfterResponseAspect = f
	return c
}
func (c *CallMethod) SetCallFunc(f func(interface{}) (*Response, error)) Caller {
	c.callFunc = f
	return c
}

func (c *CallMethod) Call(req interface{}) (*Response, error) {
	// Check if callFunc is set.
	if c.callFunc == nil {
		return nil, errors.New("callFunc is not set")
	}
	//TODO
	ctx := context.Background()
	// 应前的切面函数
	if c.BeforeRequestAspect != nil {
		err := c.BeforeRequestAspect(ctx, c.Url, c.Method, req)
		if err != nil {
			return nil, err
		}
	}

	// 执行调用
	rsp, err := c.callFunc(req)
	if err != nil {
		return nil, err
	}

	// 应后的切面函数
	if c.AfterResponseAspect != nil {
		err := c.AfterResponseAspect(ctx, rsp.Body, rsp.Status)
		if err != nil {
			return nil, err
		}
	}

	return rsp, nil
}

type Response struct {
	Body   []byte
	Status int
}
