package events

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// SetHeaders sets the required headers for SSE connection
func SetHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
}

// SendMessage sends an SSE message to the client
func SendMessage(c *gin.Context, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		return err
	}

	c.Writer.Flush()
	return nil
}

// SendEventWithType sends an SSE message with a specific event type
func SendEventWithType(c *gin.Context, eventType string, data interface{}) error {
	event := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: eventType,
		Data: data,
	}

	return SendMessage(c, event)
}
