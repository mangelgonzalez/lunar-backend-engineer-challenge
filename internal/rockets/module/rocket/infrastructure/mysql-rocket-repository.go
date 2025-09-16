package infrastructure

import (
	"database/sql"
	"errors"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/logger"
	"lunar-backend-engineer-challenge/pkg/mysql"
	"strings"

	sq "github.com/Masterminds/squirrel"
	ext_mysql "github.com/go-sql-driver/mysql"
)

type RocketRepositoryMysql struct {
	connectionPool *mysql.ConnectionPool
	tableName      string
	fields         []string
	logger         logger.Logger
}

func NewMysqlRocketRepository(connectionPool *mysql.ConnectionPool, tableName string, logger logger.Logger) *RocketRepositoryMysql {
	return &RocketRepositoryMysql{
		connectionPool: connectionPool,
		tableName:      tableName,
		fields: []string{
			"id",
			"class",
			"launch_speed",
			"mission",
		},
		logger: logger,
	}
}

func (r *RocketRepositoryMysql) scan(err error, res *sql.Rows) (*domain.Rocket, error) {
	var (
		id          string
		class       string
		launchSpeed uint
		mission     string
	)
	err = res.Scan(&id, &class, &launchSpeed, &mission)
	if err != nil {
		return nil, err
	}

	return domain.NewRocketFromPrimitives(
		id,
		class,
		launchSpeed,
		mission,
	), nil
}

func (r *RocketRepositoryMysql) FindById(id domain.RocketId) (*domain.Rocket, error) {
	queryBuilder := sq.Select(strings.Join(r.fields, ",")).
		From(r.tableName).
		Where(sq.Eq{"id": string(id)})

	res, err := queryBuilder.RunWith(r.connectionPool.Reader()).Query()

	defer mysql.CloseRows(res)

	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, domain.RocketNotExistsFromId(id)
	}

	return r.scan(err, res)
}

func (r *RocketRepositoryMysql) FindAll() (domain.Rockets, error) {
	queryBuilder := sq.Select(strings.Join(r.fields, ",")).
		From(r.tableName).
		OrderBy("class ASC")

	res, err := queryBuilder.RunWith(r.connectionPool.Reader()).Query()

	defer mysql.CloseRows(res)

	if err != nil {
		return nil, err
	}

	rockets := domain.Rockets{}
	for res.Next() {
		rocket, err := r.scan(res.Err(), res)
		if err != nil {
			return nil, err
		}

		rockets = append(rockets, rocket)
	}

	return rockets, nil
}

func (r *RocketRepositoryMysql) Save(rocket *domain.Rocket) error {
	var mysqlErr *ext_mysql.MySQLError

	rows, err := sq.Replace(r.tableName).Columns(
		r.fields...,
	).Values(
		rocket.Id(),
		rocket.Class(),
		rocket.LaunchSpeed(),
		rocket.Mission(),
	).RunWith(r.connectionPool.Writer()).Query()

	defer mysql.CloseRows(rows)

	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err = domain.NewRocketAlreadyExists(rocket.Id())
		}

		return err
	}

	return err
}

func (r *RocketRepositoryMysql) DeleteById(id domain.RocketId) error {
	queryBuilder := sq.Delete(r.tableName).
		Where(sq.Eq{"id": string(id)})

	res, err := queryBuilder.RunWith(r.connectionPool.Writer()).Query()
	defer mysql.CloseRows(res)

	return err
}
