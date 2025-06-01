package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Methods []string
type Func func(http.ResponseWriter, *http.Request)

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

func (h *Handler) Middleware(name string, mwf ...mux.MiddlewareFunc) {
	r := h.router.NewRoute().Subrouter()
	r.Use(mwf...)
	h.mws[name] = r
}

func (h *Handler) Handle(params Route) {
	r := h.router

	if m, ok := h.mws[params.Middleware]; ok {
		r = m
	}

	route := r.HandleFunc(params.Path, params.Func)
	if params.Methods != nil {
		route.Methods(params.Methods...)
	}
}

func (h *Handler) HandleRoutes(routes ...Route) {
	for _, route := range routes {
		h.Handle(route)
	}
}
