package bullhorn

// RequestResponse models subscription request response
type RequestResponse struct {
	RequestId int     `json:"requestId"`
	Events    []Event `json:"events"`
}

// Event models events from subscription request response
type Event struct {
	EventId           string   `json:"eventId"`
	EventTimestamp    int64    `json:"eventTimestamp"`
	EntityName        string   `json:"entityName"`
	EntityId          int      `json:"entityId"`
	EntityEventType   string   `json:"entityEventType"`
	UpdatedProperties []string `json:"updatedProperties"`
}

type Entity struct {
	Id        int
	Name      string
	Changes   string
	Timestamp int64
}
