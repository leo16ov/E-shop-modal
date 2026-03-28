package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	RWriter http.ResponseWriter
	Request *http.Request
	Ctx     context.Context
	values  map[string]interface{}
}

func (c *Context) Send(text string) {
	c.RWriter.Write([]byte(text))
}
func (c *Context) Status(code int) {
	c.RWriter.WriteHeader(code)
}
func (c *Context) JSONResponse(code int, data interface{}) error {
	c.RWriter.Header().Set("Content-Type", "application/json")
	c.RWriter.WriteHeader(code)
	return json.NewEncoder(c.RWriter).Encode(data)
}
func (c *Context) BindJSON(dest interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(dest)
}

func (c *Context) Set(key string, value interface{}) {
	if c.values == nil {
		c.values = make(map[string]interface{})
	}
	c.values[key] = value
}
func (c *Context) Get(key string) interface{} {
	return c.values[key]
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) Context() context.Context {
	return c.Ctx
}
