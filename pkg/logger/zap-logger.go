package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapContextLoggerOptions struct {
	zapOptions        []zap.Option
	contextDecorators ContextDecorators
}

type ZapContextLoggerOptionsFn func(options *zapContextLoggerOptions)

func WithZapOptions(zapOpts ...zap.Option) ZapContextLoggerOptionsFn {
	return func(options *zapContextLoggerOptions) {
		options.zapOptions = zapOpts
	}
}

func WithContextDecorator(decorator ContextDecorator) ZapContextLoggerOptionsFn {
	return func(options *zapContextLoggerOptions) {
		options.contextDecorators = append(options.contextDecorators, decorator)
	}
}

type ZapLoggerConfigOpt func(*zap.Config)

func SetDebugLevel(config *zap.Config) {
	config.Level.SetLevel(zap.DebugLevel)
}
func SetInfoLevel(config *zap.Config) {
	config.Level.SetLevel(zap.InfoLevel)
}
func SetWarnLevel(config *zap.Config) {
	config.Level.SetLevel(zap.WarnLevel)
}
func SetErrorLevel(config *zap.Config) {
	config.Level.SetLevel(zap.ErrorLevel)
}
func WithDateInMillisAsTime(config *zap.Config) {
	config.EncoderConfig.TimeKey = "date"
	config.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
}
func NewZapLoggerConfig(opts ...ZapLoggerConfigOpt) zap.Config {
	loggerConfig := zap.NewProductionConfig()
	SetWarnLevel(&loggerConfig)
	for _, opt := range opts {
		opt(&loggerConfig)
	}
	return loggerConfig
}

type ZapContextLogger struct {
	zap.Logger
	decorators ContextDecorators
}

func NewZapContextLogger(config zap.Config, opts ...ZapContextLoggerOptionsFn) *ZapContextLogger {
	options := &zapContextLoggerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	logger, err := config.Build(options.zapOptions...)
	if err != nil {
		panic("error building zap logger")
	}

	return &ZapContextLogger{
		Logger:     *logger,
		decorators: options.contextDecorators,
	}
}

func (zl ZapContextLogger) WithContextDecorators(decorators ...ContextDecorator) *ZapContextLogger {
	return &ZapContextLogger{
		Logger:     zl.Logger,
		decorators: append(zl.decorators, decorators...),
	}
}

func (zl ZapContextLogger) FatalContext(ctx context.Context, message string, fields ...zap.Field) {
	ctxFields := zl.decorators.Decorate(ctx)
	zl.Logger.Fatal(message, append(ctxFields, zl.errToContext(fields)...)...)
}

func (zl ZapContextLogger) ErrorContext(ctx context.Context, message string, fields ...zap.Field) {
	ctxFields := zl.decorators.Decorate(ctx)
	zl.Logger.Error(message, append(ctxFields, zl.errToContext(fields)...)...)
}

func (zl ZapContextLogger) WarnContext(ctx context.Context, message string, fields ...zap.Field) {
	ctxFields := zl.decorators.Decorate(ctx)
	zl.Logger.Warn(message, append(ctxFields, zl.errToContext(fields)...)...)
}

func (zl ZapContextLogger) InfoContext(ctx context.Context, message string, fields ...zap.Field) {
	ctxFields := zl.decorators.Decorate(ctx)
	zl.Logger.Info(message, append(ctxFields, zl.errToContext(fields)...)...)
}

func (zl ZapContextLogger) errToContext(fields []zap.Field) []zap.Field {
	newFields := []zap.Field{}
	for _, field := range fields {
		if field.Type == zapcore.ErrorType {
			newFields = append(newFields, errorToContext(field.Interface.(error))...)
			continue
		}
		newFields = append(newFields, field)
	}

	return newFields
}
