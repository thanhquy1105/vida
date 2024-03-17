package controller

import "fmt"

const (
	commonError = "ERROR"
	clientError = "CLIENT_ERROR"
)

// Error is an error that has a type
type Error struct {
	Type string
	Msg  string
}

// NewError returns a new error
func NewError(kind string, err error) *Error {
	return &Error{kind, err.Error()}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s %s", err.Type, err.Msg)
}

var (
	// ErrUnknownCommand is returned when command was not recognized
	ErrUnknownCommand = &Error{commonError, "Unknown command"}

	// ErrInvalidCommand means command wasn't parsed correcty
	ErrInvalidCommand = &Error{clientError, "Invalid command"}

	// ErrCloseCurrentItemFirst is returned when client attemted
	// to read next item before closing the current one
	ErrCloseCurrentItemFirst = &Error{clientError, "Close current item first"}

	// ErrBadDataChunk is returned when data provided by client has different size
	ErrBadDataChunk = &Error{clientError, "bad data chunk"}

	// ErrInvalidDataSize is returned when data size field is not a number
	ErrInvalidDataSize = &Error{clientError, "Invalid <bytes> number"}

	// ErrClientQuit is returned when client sends 'quit' command (not an error)
	ErrClientQuit = &Error{commonError, "Quit command received"}
)
