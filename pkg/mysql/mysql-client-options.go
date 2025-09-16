package mysql

const (
	DefaultTracer TracerType = "none"
)

type TracerType string

type Credentials struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Database string
}

func NewCredentials(user string, password string, host string, port uint16, database string) Credentials {
	return Credentials{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

type ClientOptionsFunc func(co *ClientOptions)

type ClientOptions struct {
	Credentials      Credentials
	DataSourceParams []string
	Tracer           TracerType
	TracerOptions    *ClientTracerOptions
	MaxConnections   int
	ConnIdle         int
	MaxLifetime      int
}

type ClientTracerOptions struct {
	TracerOptions []interface{}
}

func NewDefaultClientOptions(credentials Credentials, opts ...ClientOptionsFunc) *ClientOptions {
	options := &ClientOptions{
		Credentials:      credentials,
		DataSourceParams: []string{"parseTime=true"},
		Tracer:           DefaultTracer,
		TracerOptions:    nil,
		MaxConnections:   0,
		ConnIdle:         2,
		MaxLifetime:      3,
	}

	options.apply(opts...)
	return options
}

func (co *ClientOptions) apply(options ...ClientOptionsFunc) *ClientOptions {
	for _, opt := range options {
		opt(co)
	}
	return co
}

func WithMaxConnections(maxConnections int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.MaxConnections = maxConnections
	}
}

func WithConnIdle(connIdle int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.ConnIdle = connIdle
	}
}

func WithMaxLifetime(maxLifetime int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.MaxLifetime = maxLifetime
	}
}

func WithWriterMode() ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.DataSourceParams = append(co.DataSourceParams, "rejectReadOnly=true")
	}
}
