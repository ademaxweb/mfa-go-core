package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Methods []string
type Func func(http.ResponseWriter, *http.Request)

type Interface interface {
	Handle(...Route)
	Middleware(string, ...mux.MiddlewareFunc)
}

type Handler struct {
	router *mux.Router
	mws    map[string]*mux.Router
}
type Route struct {
	Path       string
	Func       Func
	Middleware string
	Methods    Methods
}

func New() *Handler {
	r := mux.NewRouter()

	return &Handler{
		router: r,
		mws:    make(map[string]*mux.Router),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) BaseMiddleware(mwf ...mux.MiddlewareFunc) {
	h.router.Use(mwf...)
}

func (h *Handler) Middleware(name string, mwf ...mux.MiddlewareFunc) {
	r := h.router.NewRoute().Subrouter()
	r.Use(mwf...)
	h.mws[name] = r
}

func (h *Handler) Handle(params ...Route) {
	r := h.router

	for _, p := range params {
		if m, ok := h.mws[p.Middleware]; ok && p.Middleware != "" {
			r = m
		}

		route := r.HandleFunc(p.Path, p.Func)
		if p.Methods != nil {
			route.Methods(p.Methods...)
		}
	}

}

func (h *Handler) HealthRoute() {
	h.Handle(Route{
		Path: "/health",
		Func: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		},
	})
}
