package migrations

type OptionsFn func(config *Options)
type Options struct {
	Platform string
}

func (ro *Options) apply(options ...OptionsFn) *Options {
	if options == nil {
		options = append(options, WithMysql())
	}
	for _, opt := range options {
		opt(ro)
	}
	return ro
}

// WithPostgres sets the database migrator to use postgres dialect.
func WithPostgres() OptionsFn {
	return func(o *Options) {
		o.Platform = postgresSqlDialect
	}
}

// WithMysql sets the database migrator to use mysql dialect.
func WithMysql() OptionsFn {
	return func(o *Options) {
		o.Platform = mysqlDialect
	}
}
