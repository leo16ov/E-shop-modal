package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	RWriter http.ResponseWriter
	Request http.Request
	Ctx     context.Context
	userID  uint
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
func (c *Context) SetUserID(id uint) {
	c.userID = id
}
func (c *Context) GetUserID() uint {
	return c.userID
}
func (c *Context) Context() context.Context {
	return c.Ctx
}
