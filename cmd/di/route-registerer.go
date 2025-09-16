package di

import (
	"lunar-backend-engineer-challenge/pkg/config"
	"lunar-backend-engineer-challenge/pkg/router"
)

type RouteRegisterer struct {
	router           *router.Router
	routesToRegister RouterFuncs
}

func NewRouteRegisterer(router *router.Router) *RouteRegisterer {
	return &RouteRegisterer{
		router: router,
	}
}

type RouterFuncs []RouterFunc
type RouterFunc func(services *CommonServices, cnf config.Config)

func (rr *RouteRegisterer) RegisterRoutes(registerer RouterFunc) {
	rr.routesToRegister = append(rr.routesToRegister, registerer)
}

func (rr *RouteRegisterer) AddRoutes(services *CommonServices, cnf config.Config) {
	for _, routerFunc := range rr.routesToRegister {
		routerFunc(services, cnf)
	}
}
