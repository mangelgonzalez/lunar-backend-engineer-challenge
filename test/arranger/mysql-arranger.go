package test

import (
	"context"
	"fmt"
	"lunar-backend-engineer-challenge/cmd/di"
	"lunar-backend-engineer-challenge/pkg/mysql"
	"sync"
)

type MysqlArranger struct {
	common *di.RocketsDI
}

func NewMysqlArranger(common *di.RocketsDI) *MysqlArranger {
	return &MysqlArranger{common: common}
}

func (sql *MysqlArranger) Arrange(_ context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := sql.common.Services.DatabaseMigrator.Up(); err != nil {
		panic(err)
	}
	sql.arrangeDB(wg)
}

func (sql *MysqlArranger) arrangeDB(wg *sync.WaitGroup) {
	rows, err := sql.common.Services.DBConnectionPool.Writer().Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}

	defer mysql.CloseRows(rows)

	var tableName string

	for rows.Next() {
		if err := rows.Scan(&tableName); nil != err {
			panic(err)
		}
		if tableName == "migrations" {
			continue
		}

		if _, err := sql.common.Services.DBConnectionPool.Writer().Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)); nil != err {
			panic(err)
		}
	}

	wg.Done()
}
