package server

import (
	"context"
	"sync"
)

type Context struct {
	context.Context
	*sync.RWMutex
}

type ContextKey struct{}

var (
	ContextRequestUUID = &ContextKey{}
)

func NewContext(uuid string) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, ContextRequestUUID, uuid)

	return &Context{
		Context: ctx,
		RWMutex: &sync.RWMutex{},
	}, cancel
}

func (c *Context) GetRequestUUID() string {
	addr, _ := c.Context.Value(ContextRequestUUID).(string)
	return addr
}

func (c *Context) SetValue(key *ContextKey, value interface{}) {
	c.Context = context.WithValue(c.Context, key, value)
}
