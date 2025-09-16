package migrations

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

type DatabaseMigrator interface {
	Up() (int, error)
	Down() (int, error)
}

const (
	postgresSqlDialect = "postgres"
	mysqlDialect       = "mysql"
)

type MysqlDatabaseMigrator struct {
	dbClient        *sql.DB
	migrationSet    *migrate.MigrationSet
	migrationSource migrate.MigrationSource
	options         *Options
}

func NewMysqlDatabaseMigrator(dbClient *sql.DB, migrationsLocation string, tableName string, opts ...OptionsFn) *MysqlDatabaseMigrator {
	source := &migrate.FileMigrationSource{
		Dir: migrationsLocation,
	}

	options := &Options{}
	options.apply(opts...)

	migrationSet := &migrate.MigrationSet{TableName: tableName}

	return &MysqlDatabaseMigrator{dbClient: dbClient, migrationSet: migrationSet, migrationSource: source, options: options}
}

func (m *MysqlDatabaseMigrator) Up() (int, error) {
	return m.migrationSet.Exec(m.dbClient, m.options.Platform, m.migrationSource, migrate.Up)
}

func (m *MysqlDatabaseMigrator) Down() (int, error) {
	return m.migrationSet.Exec(m.dbClient, m.options.Platform, m.migrationSource, migrate.Down)
}
