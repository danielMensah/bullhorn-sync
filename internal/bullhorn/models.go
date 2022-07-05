package bullhorn

// EventType is the type of event emitted by Bullhorn
type EventType string

const (
	// EventTypeInserted is the type of event emitted when an entity is inserted
	EventTypeInserted EventType = "INSERTED"
	// EventTypeUpdated is the type of event emitted when an entity is updated
	EventTypeUpdated EventType = "UPDATED"
	// EventTypeDeleted is the type of event emitted when an entity is deleted
	EventTypeDeleted EventType = "DELETED"
)

// RequestResponse models subscription request response
type RequestResponse struct {
	RequestId int     `json:"requestId"`
	Events    []Event `json:"events"`
}

// Event models events from subscription request response
type Event struct {
	EventId           string    `json:"eventId"`
	EventTimestamp    int64     `json:"eventTimestamp"`
	EntityName        string    `json:"entityName"`
	EntityId          int32     `json:"entityId"`
	EntityEventType   EventType `json:"entityEventType"`
	UpdatedProperties []string  `json:"updatedProperties"`
}

// Entity models entities from subscription request response
type Entity struct {
	Id        int32
	Name      string
	EventType EventType
	Changes   []byte
	Timestamp int64
}
