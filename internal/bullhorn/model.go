package bullhorn

import "time"

const (
	inserted = "INSERTED"
	updated  = "UPDATED"
	deleted  = "DELETED"
)

type RequestResponse struct {
	RequestId int     `json:"requestId"`
	Events    []Event `json:"events"`
}

type Event struct {
	EventId           string   `json:"eventId"`
	EventTimestamp    int64    `json:"eventTimestamp"`
	EntityName        string   `json:"entityName"`
	EntityId          int      `json:"entityId"`
	EntityEventType   string   `json:"entityEventType"`
	UpdatedProperties []string `json:"updatedProperties"`
}

type EventMetadata struct {
	PersonID      string `json:"PERSON_ID"`
	TransactionID string `json:"TRANSACTION_ID"`
}

type Record struct {
	EntityId        int
	EntityName      string
	EntityEventType string
	EventTimestamp  time.Time
	Changes         []byte
}
