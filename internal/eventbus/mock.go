package eventbus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/stretchr/testify/mock"
)

// Mock mocks eventBridge methods for use in the services using the eventBridge client
type Mock struct {
	eventbridgeiface.EventBridgeAPI
	mock.Mock
}

// PutEventsWithContext mocks this method used for testing
func (m *Mock) PutEventsWithContext(ctx aws.Context, input *eventbridge.PutEventsInput, opts ...request.Option) (*eventbridge.PutEventsOutput, error) {
	args := m.Called(ctx, input)
	eventsOutput, ok := args.Get(0).(*eventbridge.PutEventsOutput)
	if !ok {
		eventsOutput = nil
	}

	return eventsOutput, args.Error(1)
}
