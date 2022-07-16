package simple_server

import (
	"net/http"
	"time"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

func (c *Context) Param(key string) string {
	return c.R.FormValue(key)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.R.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.R.Context().Done()
}

func (c *Context) Err() error {
	return c.R.Context().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.R.Context().Value(key)
}
