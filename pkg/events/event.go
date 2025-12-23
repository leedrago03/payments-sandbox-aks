package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents a system event.
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// NewEvent creates a new Event.
func NewEvent(eventType, source string, data interface{}) (*Event, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now(),
		Data:      jsonData,
		Metadata:  make(map[string]string),
	}, nil
}
