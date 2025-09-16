package command

import (
	"context"
	"lunar-backend-engineer-challenge/pkg/bus"
	"lunar-backend-engineer-challenge/pkg/logger"
	mutex "lunar-backend-engineer-challenge/pkg/sync"
	"reflect"
	"sync"

	"go.uber.org/zap"
)

type Bus interface {
	RegisterCommand(command bus.Dto, handler CommandHandler) error
	GetHandler(command bus.Dto) (CommandHandler, error)
	Dispatch(ctx context.Context, dto bus.Dto) error
	DispatchAsync(ctx context.Context, dto bus.Dto) error
	ProcessFailed(ctx context.Context)
}

type CommandBus struct {
	handlers       map[string]CommandHandler
	lock           sync.Mutex
	logger         logger.Logger
	failedCommands chan *FailedCommand

	mutex mutex.MutexService
}

func InitCommandBus(logger logger.Logger, mutex mutex.MutexService) *CommandBus {
	return &CommandBus{
		handlers:       make(map[string]CommandHandler, 0),
		lock:           sync.Mutex{},
		logger:         logger,
		failedCommands: make(chan *FailedCommand),

		mutex: mutex,
	}
}

type FailedCommand struct {
	command        bus.Dto
	handler        CommandHandler
	timesProcessed int
}

type CommandAlreadyRegistered struct {
	message     string
	commandName string
}

func (i CommandAlreadyRegistered) Error() string {
	return i.message
}

func NewCommandAlreadyRegistered(message string, commandName string) CommandAlreadyRegistered {
	return CommandAlreadyRegistered{message: message, commandName: commandName}
}

type CommandNotRegistered struct {
	message     string
	commandName string
}

func (i CommandNotRegistered) Error() string {
	return i.message
}

func NewCommandNotRegistered(message string, commandName string) CommandNotRegistered {
	return CommandNotRegistered{message: message, commandName: commandName}
}

func (cb *CommandBus) RegisterCommand(command bus.Dto, handler CommandHandler) error {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	commandName, err := cb.commandName(command)
	if err != nil {
		return err
	}

	if _, ok := cb.handlers[*commandName]; ok {
		return NewCommandAlreadyRegistered("Command already registered", *commandName)
	}

	cb.handlers[*commandName] = handler

	return nil
}

func (cb *CommandBus) GetHandler(command bus.Dto) (CommandHandler, error) {
	commandName, err := cb.commandName(command)
	if err != nil {
		return nil, err
	}
	if handler, ok := cb.handlers[*commandName]; ok {
		return handler, nil
	}

	return nil, NewCommandNotRegistered("Command not registered", *commandName)
}

func (cb *CommandBus) Dispatch(ctx context.Context, command bus.Dto) error {
	handler, err := cb.GetHandler(command)
	if err != nil {
		return err
	}

	return cb.doHandle(ctx, handler, command)
}

func (cb *CommandBus) DispatchAsync(ctx context.Context, command bus.Dto) error {
	commandName, err := cb.commandName(command)
	if err != nil {
		return err
	}

	if handler, ok := cb.handlers[*commandName]; ok {
		go cb.doHandleAsync(ctx, handler, command)

		return nil
	}

	return NewCommandNotRegistered("Command not registered", *commandName)
}

func (cb *CommandBus) doHandle(ctx context.Context, handler CommandHandler, command bus.Dto) error {

	if bc, ok := command.(bus.BlockOperationCommand); ok {
		operation := func() (interface{}, error) {
			return nil, handler.Handle(ctx, bc)
		}

		_, err := cb.mutex.Mutex(ctx, bc.BlockingKey(), operation)

		return err
	}

	return handler.Handle(ctx, command)
}

func (cb *CommandBus) doHandleAsync(ctx context.Context, handler CommandHandler, command bus.Dto) {
	err := cb.doHandle(ctx, handler, command)

	if err != nil {
		cb.failedCommands <- &FailedCommand{
			command:        command,
			handler:        handler,
			timesProcessed: 1,
		}
		cb.logger.Error(err.Error())
	}
}

func (cb *CommandBus) commandName(cmd interface{}) (*string, error) {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return nil, CommandNotValid{"only pointer to commands are allowed"}
	}

	name := value.String()

	return &name, nil
}

func (cb *CommandBus) ProcessFailed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(cb.failedCommands)
			cb.logger.Warn("Exiting safely failed commands consumer...")
			return
		case failedCommand := <-cb.failedCommands:
			if failedCommand.timesProcessed >= 3 {
				continue
			}

			failedCommand.timesProcessed++
			if err := cb.doHandle(ctx, failedCommand.handler, failedCommand.command); err != nil {
				cb.logger.Warn(err.Error(), zap.Any("previous_error", err))
			}
		}
	}
}

type CommandNotValid struct {
	message string
}

func (i CommandNotValid) Error() string {
	return i.message
}
