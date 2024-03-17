package controller

import (
	"bufio"
	"io"
	"time"

	"github.com/thanhquy1105/vida/repository"
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
