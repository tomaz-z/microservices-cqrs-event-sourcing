package commands

import (
	escqrs "eventsourcing/services"
)

type command struct {
	domain      string
	commandType string
	data        interface{}
}

// New returns a new command.
func New(
	domain string,
	commandType string,
	data interface{},
) escqrs.Command {
	return &command{
		domain:      domain,
		commandType: commandType,
		data:        data,
	}
}

func (c command) Domain() string {
	return c.domain
}

func (c command) Type() string {
	return c.commandType
}

func (c command) Data() interface{} {
	return c.data
}
