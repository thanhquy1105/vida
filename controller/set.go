package controller

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Set handles SET command
// Command: SET <queue> <not_impl> <not_impl> <bytes>
// <data block>
// Response: STORED
func (c *Controller) Set(input []string) error {
	cmd, err := parseSetCommand(input)
	if err != nil {
		return err
	}

	// After that cmd, the client sends the data block
	dataBlock, err := c.readDataBlock(cmd.DataSize)
	if err != nil {
		return err
	}

	if cmd.FanoutQueues == nil {
		err = c.storeDataBlock(cmd.QueueName, dataBlock)
		if err != nil {
			log.Println(cmd, err)
			return err
		}
	} else {
		for _, queueName := range cmd.FanoutQueues {
			err = c.storeDataBlock(queueName, dataBlock)
			if err != nil {
				log.Println(cmd, err)
				return err
			}
		}
	}

	fmt.Fprint(c.rw.Writer, storedMessage)
	c.rw.Writer.Flush()
	return nil
}

func (c *Controller) readDataBlock(totalBytes int) ([]byte, error) {
	// makes new buffer for larger data block
	// or use the same one
	if cap(c.dataBuffer) < totalBytes+2 {
		c.dataBuffer = make([]byte, totalBytes+2)
	} else {
		c.dataBuffer = c.dataBuffer[:totalBytes+2]
	}

	// <data block>\r\n
	_, err := io.ReadFull(c.rw.Reader, c.dataBuffer)
	if err != nil {
		return nil, err
	}

	if string(c.dataBuffer[totalBytes:]) != "\r\n" {
		return nil, ErrBadDataChunk
	}

	return c.dataBuffer[:totalBytes], nil
}

func (c *Controller) storeDataBlock(queueName string, dataBlock []byte) error {
	fmt.Println(queueName)
	fmt.Println(dataBlock)
	return nil
}

func parseSetCommand(input []string) (*Command, error) {
	if len(input) < 5 || len(input) > 6 {
		return nil, ErrInvalidCommand
	}

	// get bytes
	totalBytes, err := strconv.Atoi(input[4])
	if err != nil {
		return nil, ErrInvalidDataSize
	}

	cmd := &Command{Name: input[0], QueueName: input[1], DataSize: totalBytes}

	// set new message into multiple queues at once by using the following syntax
	// set <queue>+<another_queue>+<third_queue> => fanout
	if strings.Contains(cmd.QueueName, "+") {
		cmd.FanoutQueues = strings.Split(cmd.QueueName, "+")
		cmd.QueueName = cmd.FanoutQueues[0]
	}
	return cmd, nil
}
