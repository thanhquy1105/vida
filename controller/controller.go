package controller

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/thanhquy1105/vida/repository"
)

const (
	storedMessage = "STORED\r\n"
	endMessage    = "END\r\n"
)

// Conn represents a connecting consumer interface
type Conn interface {
	io.Reader
	io.Writer
	SetDeadline(t time.Time) error
}

// Controller represents a controller of connecting consumer
type Controller struct {
	conn       Conn
	rw         *bufio.ReadWriter
	repo       *repository.QueueRepository
	dataBuffer []byte
}

// Command represents a comsumer command
type Command struct {
	Name          string
	QueueName     string
	ConsumerGroup string
	// FanoutQueues is the queue name array to save value into
	FanoutQueues []string
	// DataSize is the data size of the value
	DataSize int
}

// NewSession creates new controller of connecting consumer
func NewSession(conn Conn, repo *repository.QueueRepository) *Controller {
	return &Controller{
		conn:       conn,
		rw:         bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		repo:       repo,
		dataBuffer: make([]byte, 1024),
	}
}

// ReadFirstMessage reads initial message from connection buffer
func (c *Controller) ReadFirstMessage() (string, error) {
	return c.rw.Reader.ReadString('\n')
}

// UnknownCommand reports an error
func (c *Controller) UnknownCommand() error {
	c.SendError(ErrUnknownCommand)
	return ErrUnknownCommand
}

// SendError sends an error message to the client
func (c *Controller) SendError(err error) {

	if e, ok := err.(*Error); ok {
		fmt.Fprintf(c.rw.Writer, "%s\r\n", e.Error())
	} else {
		fmt.Fprintf(c.rw.Writer, "%s %s\r\n", commonError, err.Error())
	}

	c.rw.Writer.Flush()
}
