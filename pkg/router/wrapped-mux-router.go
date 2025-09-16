package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type WrappedMuxRouter interface {
	Use(mwf ...mux.MiddlewareFunc)
	NewRoute() *mux.Route
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

func optsFromMuxRouterOptions(muxOpts MuxRouterOptions) []WrappedMuxRouterOpt {
	opts := make([]WrappedMuxRouterOpt, 0)
	if muxOpts.SkipClean {
		opts = append(opts, WithSkipClean())
	}
	if muxOpts.StrictSlash {
		opts = append(opts, WithStrictSlash())
	}
	if muxOpts.UseEncodedPath {
		opts = append(opts, WithEncodedPath())
	}
	return opts
}

func applyOpts(wrappedRouter *mux.Router, opts ...WrappedMuxRouterOpt) {
	for _, opt := range opts {
		opt(wrappedRouter)
	}
}

type WrappedMuxRouterOpt func(router *mux.Router)

func WithSkipClean() WrappedMuxRouterOpt {
	return func(router *mux.Router) {
		router.SkipClean(true)
	}
}

func WithStrictSlash() WrappedMuxRouterOpt {
	return func(router *mux.Router) {
		router.StrictSlash(true)
	}
}

func WithEncodedPath() WrappedMuxRouterOpt {
	return func(router *mux.Router) {
		router.UseEncodedPath()
	}
}
