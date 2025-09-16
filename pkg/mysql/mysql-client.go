package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const driverName = "mysql"

func Dial(user string, password string, host string, port uint16, databaseName string) (*sql.DB, error) {

	sqlAddress := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, databaseName)

	dbClient, err := sql.Open(driverName, sqlAddress)
	if err != nil {
		return nil, err
	}

	if err = dbClient.Ping(); err != nil {
		return nil, err
	}

	dbClient.SetConnMaxLifetime(3 * time.Minute)

	return dbClient, nil
}

func CloseRows(rows *sql.Rows) {
	if rows != nil {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		if err := stmt.Close(); err != nil {
			panic(err)
		}
	}
}

func NewReader(
	user string,
	password string,
	host string,
	port uint16,
	databaseName string,
	maxConnections int,
	connIdle int,
	maxLifetimeInMinutes int,
) (*sql.DB, error) {
	credentials := NewCredentials(user, password, host, port, databaseName)

	return NewClientReader(credentials, WithMaxConnections(maxConnections), WithConnIdle(connIdle), WithMaxLifetime(maxLifetimeInMinutes))
}

func NewWriter(
	user string,
	password string,
	host string,
	port uint16,
	databaseName string,
	maxConnections int,
	connIdle int,
	maxLifetimeInMinutes int,
) (*sql.DB, error) {
	credentials := NewCredentials(user, password, host, port, databaseName)

	return NewClientWriter(credentials, WithMaxConnections(maxConnections), WithConnIdle(connIdle), WithMaxLifetime(maxLifetimeInMinutes))
}

func NewClientReader(credentials Credentials, opt ...ClientOptionsFunc) (*sql.DB, error) {
	options := NewDefaultClientOptions(credentials, opt...)
	return newClient(options)
}

func NewClientWriter(credentials Credentials, opt ...ClientOptionsFunc) (*sql.DB, error) {
	opt = append(opt, WithWriterMode())
	options := NewDefaultClientOptions(credentials, opt...)

	return newClient(options)
}

func newClient(ops *ClientOptions) (*sql.DB, error) {
	var dbClient *sql.DB
	var err error

	sqlAddress := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		ops.Credentials.User,
		ops.Credentials.Password,
		ops.Credentials.Host,
		ops.Credentials.Port,
		ops.Credentials.Database,
		strings.Join(ops.DataSourceParams, "&"),
	)

	switch ops.Tracer {
	default:
		dbClient, err = defaultConnection(sqlAddress)
	}

	if err != nil {
		return dbClient, err
	}

	if err = dbClient.Ping(); err != nil {
		return nil, err
	}

	dbClient.SetConnMaxLifetime(time.Duration(ops.MaxLifetime) * time.Minute)
	dbClient.SetMaxOpenConns(ops.MaxConnections)
	dbClient.SetMaxIdleConns(ops.ConnIdle)

	return dbClient, nil
}

func defaultConnection(sqlAddress string) (*sql.DB, error) {
	return sql.Open(driverName, sqlAddress)
}
