package server

import "net/http"

func (app *App) Get(path string, handler func(*Context)) {
	app.mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	app.handlerCount++
}
func (app *App) Post(path string, handler func(*Context)) {
	app.mux.HandleFunc("POST "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	app.handlerCount++
}
func (app *App) Put(path string, handler func(*Context)) {
	app.mux.HandleFunc("PUT "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	app.handlerCount++
}
func (app *App) Delete(path string, handler func(*Context)) {
	app.mux.HandleFunc("DELETE "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	app.handlerCount++
}
