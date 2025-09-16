package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type Router struct {
	negroni *negroni.Negroni
	mux     WrappedMuxRouter
	server  *http.Server
}

type Options struct {
	WriteTimeout     int
	ReadTimeout      int
	CorsOptions      *cors.Options
	Middlewares      []Middleware
	RoutedMiddleware []Middleware
	MuxOptions       MuxRouterOptions
}

type DatadogTracerOptions struct {
	Active      bool
	ServiceName string
}

type MuxRouterOptions struct {
	StrictSlash    bool
	SkipClean      bool
	UseEncodedPath bool
}

const HeaderVersionName = "Accept-version"

type RouteMatching func(r *http.Request) bool

type Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func NewRouter(options Options) Router {
	readTimeout := 30 * time.Second
	writeTimeout := 30 * time.Second

	if options.ReadTimeout != 0 {
		readTimeout = time.Duration(options.ReadTimeout) * time.Second
	}

	if options.WriteTimeout != 0 {
		writeTimeout = time.Duration(options.WriteTimeout) * time.Second
	}

	m := muxRouter(options.MuxOptions)
	n := negroni.New()

	if options.CorsOptions != nil {
		c := cors.New(*options.CorsOptions)

		n.Use(c)
	}

	for _, middleware := range options.Middlewares {
		n.UseFunc(middleware)
	}

	h := &http.Server{
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	return Router{n, m, h}
}

func muxRouter(muxOptions MuxRouterOptions) WrappedMuxRouter {
	rr := mux.NewRouter()
	applyOpts(rr, optsFromMuxRouterOptions(muxOptions)...)

	return rr
}

func parameterizedRouter(serverWriteTimeOut int, serverReadTimeout int, muxOptions MuxRouterOptions, middlewares ...Middleware) *Router {
	options := Options{
		WriteTimeout: serverReadTimeout,
		ReadTimeout:  serverWriteTimeOut,
		CorsOptions: &cors.Options{
			AllowedOrigins:   []string{"*", "http://localhost*"},
			AllowedMethods:   []string{"POST", "GET", "HEAD", "PATCH", "OPTIONS", "PUT"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		},
		Middlewares: middlewares,
		MuxOptions:  muxOptions,
	}

	rt := NewRouter(options)

	return &rt
}

func DefaultRouter(serverWriteTimeOut int, serverReadTimeout int, middlewares ...Middleware) *Router {
	return parameterizedRouter(serverWriteTimeOut, serverReadTimeout, MuxRouterOptions{StrictSlash: true}, middlewares...)
}

func (router *Router) ListenAndServe(addr string) error {
	router.server.Addr = addr
	router.server.Handler = router.BuildHandler()

	// Use default options
	return router.server.ListenAndServe()
}

func (router *Router) BuildHandler() http.Handler {
	router.negroni.UseHandler(router.mux)
	return router.negroni
}

func (router *Router) Handle(method string, path string, handler http.Handler, routeMatching RouteMatching, middlewares ...Middleware) {
	router.handleMultipleMethods([]string{method}, path, handler, routeMatching, middlewares...)
}

func (router *Router) handleMultipleMethods(methods []string, path string, handler http.Handler, routeMatching RouteMatching, middlewares ...Middleware) {
	var stack []negroni.Handler

	for _, middleware := range middlewares {
		stack = append(stack, negroni.HandlerFunc(middleware))
	}

	stack = append(stack, negroni.Wrap(handler))

	route := router.mux.NewRoute()

	route.Handler(negroni.New(stack...))

	if nil != routeMatching {
		route.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return routeMatching(r)
		})
	}

	if strings.HasSuffix(path, "*") {
		route.PathPrefix(path[:len(path)-1])
	} else {
		route.Path(path)
	}

	if len(methods) != 0 {
		route.Methods(methods...)
	}
}

func (router *Router) Get(path string, routeMatching RouteMatching, handler http.Handler, middlewares ...Middleware) {
	router.Handle("GET", path, handler, routeMatching, middlewares...)
}

func (router *Router) Post(path string, routeMatching RouteMatching, handler http.Handler, middlewares ...Middleware) {
	router.Handle("POST", path, handler, routeMatching, middlewares...)
}

func NewDefaultRouteMatching() RouteMatching {
	return NewHeaderRouteMatching(HeaderVersionName, "")
}

func NewHeaderRouteMatching(headerName string, headerValue string) RouteMatching {
	return func(r *http.Request) bool {
		return r.Header.Get(headerName) == headerValue
	}
}
