package server

import "net/http"

func HandleFunc(
	mux *http.ServeMux,
	pattern string,
	handler func(*Context),
) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		}
		handler(ctx)
	})
}
func HandleProtected(
	mux *http.ServeMux,
	pattern string,
	handler func(*Context),
	middleware func(func(*Context)) func(*Context),
) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		}

		finalHandler := middleware(handler)
		finalHandler(ctx)
	})
}
