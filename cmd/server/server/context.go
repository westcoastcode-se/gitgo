package server

import (
	"context"
	"net"
	"sync"
)

type Context struct {
	context.Context
	*sync.RWMutex
}

type ContextKey struct{}

var (
	// ContextPermissions used when fetching the permissions associated with a specific request
	ContextPermissions = &ContextKey{}
	ContextUser        = &ContextKey{}
	ContextLocalAddr   = &ContextKey{}
	ContextRemoteAddr  = &ContextKey{}
)

func NewContext() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		Context: ctx,
		RWMutex: &sync.RWMutex{},
	}, cancel
}

// GetPermissions fetches the permissions associated with a specific request
func (c *Context) GetPermissions() *Permissions {
	val, ok := c.Value(ContextPermissions).(*Permissions)
	if ok {
		return val
	}
	return MissingPermissions
}

func (c *Context) GetLocalAddr() net.Addr {
	addr, _ := c.Context.Value(ContextLocalAddr).(net.Addr)
	return addr
}

func (c *Context) GetRemoteAddr() net.Addr {
	addr, _ := c.Context.Value(ContextRemoteAddr).(net.Addr)
	return addr
}

func (c *Context) SetValue(key *ContextKey, value interface{}) {
	c.Context = context.WithValue(c.Context, key, value)
}
