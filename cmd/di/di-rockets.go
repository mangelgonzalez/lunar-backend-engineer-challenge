package di

import (
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/application"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/infrastructure"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/ui"
	"lunar-backend-engineer-challenge/pkg/bus/command"
	"lunar-backend-engineer-challenge/pkg/bus/query"
	"lunar-backend-engineer-challenge/pkg/config"
	"lunar-backend-engineer-challenge/pkg/router"
)

type RocketsServices struct {
	Repository                        domain.RocketRepository
	FindRocketByQueryHandler          *application.FindRocketByIdQueryHandler
	FindAllRocketsQueryHandler        *application.FindAllRocketsQueryHandler
	StoreRocketMessagesCommandHandler *application.StoreRocketMessagesCommandHandler
}

func InitRocketsServices(services *CommonServices) *RocketsServices {
	repository := infrastructure.NewMysqlRocketRepository(services.DBConnectionPool, "rockets", services.Logger)
	storer := domain.NewRocketStorer(repository)

	rocketsServices := &RocketsServices{
		Repository:                        repository,
		FindRocketByQueryHandler:          application.NewFindRocketByIdQueryHandler(repository),
		FindAllRocketsQueryHandler:        application.NewFindAllRocketsQueryHandlerHandler(repository),
		StoreRocketMessagesCommandHandler: application.NewStoreRocketMessagesCommandHandler(storer),
	}

	registerRocketsQueries(services.QueryBus, rocketsServices)
	registerRocketsCommands(services.CommandBus, rocketsServices)
	services.RegisterRoutes(registerRocketsRoutes())

	return rocketsServices
}

func registerRocketsQueries(bus query.Bus, service *RocketsServices) {
	if err := bus.RegisterQuery(&application.FindRocketByIdQuery{}, service.FindRocketByQueryHandler); err != nil {
		panic(err)
	}

	if err := bus.RegisterQuery(&application.FindAllRocketsQuery{}, service.FindAllRocketsQueryHandler); err != nil {
		panic(err)
	}
}

func registerRocketsCommands(bus command.Bus, service *RocketsServices) {
	if err := bus.RegisterCommand(&application.StoreRocketMessagesCommand{}, service.StoreRocketMessagesCommandHandler); err != nil {
		panic(err)
	}
}

func registerRocketsRoutes() RouterFunc {
	return func(services *CommonServices, cnf config.Config) {
		routeMatching := router.NewDefaultRouteMatching()

		services.router.Get(
			"/rockets/{rocketId}",
			routeMatching,
			ui.HandleGetRocketById(services.QueryBus, services.JsonApiResponseMiddleware),
		)

		services.router.Get(
			"/rockets",
			routeMatching,
			ui.HandleGetAllRockets(services.QueryBus, services.JsonApiResponseMiddleware),
		)

		services.router.Post(
			"/rockets",
			routeMatching,
			ui.HandleStoreRocketMessages(services.CommandBus, services.JsonApiResponseMiddleware),
		)
	}

}
