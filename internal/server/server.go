package server

import (
	"log"
	"net/http"
)

type App struct {
	mux          *http.ServeMux
	handlerCount int
}

func NewApp() *App {
	return &App{
		mux:          http.NewServeMux(),
		handlerCount: 0,
	}
}
func (app *App) RunServer(port string) error {
	server := &http.Server{
		Addr:    port,
		Handler: app.mux,
	}
	log.Printf("Run Server: http://localhost%s\n", port)
	return server.ListenAndServe()
}
