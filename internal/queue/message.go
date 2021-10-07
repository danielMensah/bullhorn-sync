package queue

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type deleter interface {
	deleteMessage(receiptHandle *string) error
}

// Message holds the notifications converted from the SQS message event as well as the sqsMessage itself with a deleter
// interface to be able to delete the event from the queue once finished processing
type Message struct {
	Event         events.CloudWatchEvent
	receiptHandle *string
	deleter       deleter
	ReceivedTime  time.Time
}

// Delete deletes messages from queue¬
func (m *Message) Delete() error {
	return m.deleter.deleteMessage(m.receiptHandle)
}
