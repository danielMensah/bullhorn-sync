package bullhorn

import (
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RequestResponse models subscription request response
type RequestResponse struct {
	RequestId int     `json:"requestId"`
	Events    []Event `json:"events"`
}

// Event models events from subscription request response
type Event struct {
	EventId           string       `json:"eventId"`
	EventTimestamp    int64        `json:"eventTimestamp"`
	EntityName        string       `json:"entityName"`
	EntityId          int32        `json:"entityId"`
	EntityEventType   pb.EventType `json:"entityEventType"`
	UpdatedProperties []string     `json:"updatedProperties"`
}

// Entity models entities from subscription request response
type Entity struct {
	Id        int32
	Name      string
	Changes   []byte
	Timestamp *timestamppb.Timestamp
}
