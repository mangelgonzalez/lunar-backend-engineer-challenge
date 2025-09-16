package di

import (
	"context"
	"database/sql"
	"fmt"
	"lunar-backend-engineer-challenge/pkg/bus/command"
	"lunar-backend-engineer-challenge/pkg/bus/query"
	"lunar-backend-engineer-challenge/pkg/config"
	"lunar-backend-engineer-challenge/pkg/http/middleware"
	"lunar-backend-engineer-challenge/pkg/logger"
	"lunar-backend-engineer-challenge/pkg/migrations"
	"lunar-backend-engineer-challenge/pkg/mysql"
	redislib "lunar-backend-engineer-challenge/pkg/redis"
	"lunar-backend-engineer-challenge/pkg/router"
	"lunar-backend-engineer-challenge/pkg/sync"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

type CommonServices struct {
	Logger logger.Logger
	Router *router.Router

	RedisClient      *redis.Client
	MutexService     sync.MutexService
	DBConnectionPool *mysql.ConnectionPool
	DatabaseMigrator migrations.DatabaseMigrator

	QueryBus                  query.Bus
	CommandBus                command.Bus
	JsonApiResponseMiddleware *middleware.JsonApiResponseMiddleware

	*RouteRegisterer
}

func Init() *RocketsDI {
	return InitWithEnvFile("./.env")
}

func InitWithEnvFile(envFiles ...string) *RocketsDI {
	err := godotenv.Overload(envFiles...)
	if nil != err {
		panic(err)
	}

	return setUp()
}

type RocketsDI struct {
	Services       *CommonServices
	Config         config.Config
	RocketServices *RocketsServices
}

func setUp() *RocketsDI {
	cnf := buildConfig()
	l := buildLogger()
	redisClient := buildRedisClient(cnf)
	redisMutexService := sync.NewRedisMutexService(redisClient, l)
	dbWriterClient := buildDBWriterClient(cnf)
	dbReaderClient := buildDBReaderClient(cnf)
	dbConnectionPool := mysql.NewMysqlConnectionPool(dbWriterClient, dbReaderClient)
	queryBus := query.InitQueryBus(l)
	commandBus := command.InitCommandBus(l, redisMutexService)
	jsonApiResponseMiddleware := middleware.NewJsonApiResponseMiddleware(l)

	rt := NewRouter(cnf)

	services := &CommonServices{
		Logger:                    l,
		Router:                    rt,
		RedisClient:               redisClient,
		MutexService:              redisMutexService,
		DBConnectionPool:          dbConnectionPool,
		DatabaseMigrator:          migrations.NewMysqlDatabaseMigrator(dbWriterClient, cnf.MigrationsLocation, cnf.MigrationsTableName, migrations.WithMysql()),
		QueryBus:                  queryBus,
		CommandBus:                commandBus,
		JsonApiResponseMiddleware: jsonApiResponseMiddleware,
		RouteRegisterer:           NewRouteRegisterer(rt),
	}

	rocketsDi := &RocketsDI{
		Services:       services,
		Config:         cnf,
		RocketServices: InitRocketsServices(services),
	}

	services.AddRoutes(services, cnf)

	return rocketsDi
}

func Context() (context.Context, context.CancelFunc) {
	rootCtx, cancel := context.WithCancel(context.Background())

	return rootCtx, cancel
}

func buildConfig() config.Config {
	var c config.Config
	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	return c
}

func buildRedisClient(config config.Config) *redis.Client {
	return redislib.InitRedisClient(
		config.RedisHost,
		config.RedisPort,
		config.RedisMaxIdleConnections,
		config.RedisIdleTimeout,
	)
}

func buildLogger() *zap.Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return l
}

func buildDBWriterClient(config config.Config) *sql.DB {
	sqlClient, err := mysql.NewWriter(
		config.MySqlUser,
		config.MySqlPass,
		config.MySqlWriterHost,
		config.MySqlPort,
		config.MySqlDB,
		0,
		2,
		3,
	)
	if err != nil {
		panic(err)
	}

	return sqlClient
}

func buildDBReaderClient(config config.Config) *sql.DB {
	sqlClient, err := mysql.NewReader(
		config.MySqlUser,
		config.MySqlPass,
		config.MySqlReaderHost,
		config.MySqlPort,
		config.MySqlDB,
		0,
		2,
		3,
	)
	if err != nil {
		panic(err)
	}

	return sqlClient
}

func RunMigrations(ctx context.Context, services *CommonServices) {
	migrationFunc := func() (interface{}, error) {
		migrationsExecuted, err := services.DatabaseMigrator.Up()
		if err != nil {
			return nil, err
		}

		return migrationsExecuted, nil
	}

	migrationsDone, err := services.MutexService.Mutex(ctx, "migrations", migrationFunc)
	if err != nil {
		panic(err)
	}

	services.Logger.Warn(fmt.Sprintf("Applied %d migrations!", migrationsDone.(int)))
}

func NewRouter(config config.Config) *router.Router {
	return router.DefaultRouter(
		config.ServerWriteTimeout,
		config.ServerReadTimeout,
	)
}
