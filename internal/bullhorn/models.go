package bullhorn

type EventType string

const (
	EventType_INSERTED EventType = "INSERT"
	EventType_UPDATED  EventType = "UPDATE"
	EventType_DELETED  EventType = "DELETE"
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
