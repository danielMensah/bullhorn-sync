package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	log "github.com/sirupsen/logrus"
)

//SQS is simply used to hold the client, event and message chan for the processing of SQS messages
type SQS struct {
	client        sqsiface.SQSAPI
	messages      chan<- Message
	receiveConfig receiveConfig
}

type receiveConfig struct {
	queueURL          string
	waitTime          int64
	visibilityTimeout int64
}

// NewSqsConsumer will create a SQS client. If client is passed through it will use the defined client
func NewSqsConsumer(region string, messages chan<- Message) (*SQS, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	})
	if err != nil {
		return nil, err
	}

	client := sqs.New(sess)

	return &SQS{
		client:   client,
		messages: messages,
		receiveConfig: receiveConfig{
			queueURL:          "sqsURL",
			waitTime:          10,
			visibilityTimeout: 120,
		},
	}, nil
}

// Receive will continuously request messages from SQS and process them accordingly.
func (s SQS) Receive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resp, err := s.client.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
					aws.String(sqs.MessageSystemAttributeNameAwstraceHeader),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            &s.receiveConfig.queueURL,
				MaxNumberOfMessages: aws.Int64(10),
				VisibilityTimeout:   &s.receiveConfig.visibilityTimeout,
				WaitTimeSeconds:     &s.receiveConfig.waitTime,
			})
			if err != nil {
				log.WithError(err).Error("receiving SQS message")
				continue
			}

			for _, message := range resp.Messages {
				receivedTime := time.Now()

				event := events.CloudWatchEvent{}
				logger := log.WithField("SQS ID", *message.MessageId)

				if err := json.Unmarshal([]byte(*message.Body), &event); err != nil {
					logger.WithError(err).Error("decoding event")

					// if we've errored decoding the message delete it from the queue
					if err := s.deleteMessage(message.ReceiptHandle); err != nil {
						logger.WithError(err).Error("deleting message during receive")
					}

					continue
				}

				s.messages <- Message{
					Event:         event,
					receiptHandle: message.ReceiptHandle,
					deleter:       s,
					ReceivedTime:  receivedTime,
				}
			}
		}
	}
}

func (s SQS) deleteMessage(receiptHandle *string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &s.receiveConfig.queueURL,
		ReceiptHandle: receiptHandle,
	}

	if _, err := s.client.DeleteMessage(input); err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	return nil
}
