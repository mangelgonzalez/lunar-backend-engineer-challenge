package mysql

import (
	"database/sql"
	"lunar-backend-engineer-challenge/pkg/domain"
)

type ConnectionPool struct {
	writer *sql.DB
	reader *sql.DB
}

func NewMysqlConnectionPool(writer *sql.DB, reader *sql.DB) *ConnectionPool {
	guardConnection(writer)
	guardConnection(reader)

	return &ConnectionPool{
		writer: writer,
		reader: reader,
	}
}

func guardConnection(conn *sql.DB) {
	if nil == conn {
		panic(NewInvalidMysqlPoolConfigProvided())
	}
}

func (c *ConnectionPool) Writer() *sql.DB {
	return c.writer
}

func (c *ConnectionPool) Reader() *sql.DB {
	if nil == c.reader {
		return c.writer
	}

	return c.reader
}

const invalidMysqlPoolConfigProvided = "Invalid MySQL pool config provided"

type InvalidMysqlPoolConfigProvided struct {
	extraItems map[string]interface{}
	domain.CriticalError
}

func (u InvalidMysqlPoolConfigProvided) Error() string {
	return invalidMysqlPoolConfigProvided
}

func (u InvalidMysqlPoolConfigProvided) ExtraItems() map[string]interface{} {
	return u.extraItems
}

func NewInvalidMysqlPoolConfigProvided() *InvalidMysqlPoolConfigProvided {
	return &InvalidMysqlPoolConfigProvided{extraItems: map[string]interface{}{}}
}
