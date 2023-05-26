package internal

import "errors"

var (
	DatabaseErr       = errors.New("database error, see log for more info")
	InvalidCommandErr = errors.New("invalid command")
)
