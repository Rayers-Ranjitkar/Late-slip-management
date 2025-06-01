package events

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEHandler handles the SSE connection
func SSEHandler(c *gin.Context) {
	userId := c.GetString("user_id")
	isAdmin := c.GetBool("is_admin")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Register client
	client := Manager.AddClient(userId, isAdmin)
	defer Manager.RemoveClient(userId)

	// Create channel for detecting client disconnection
	clientGone := c.Writer.CloseNotify()

	// Send initial connection established event
	sendEvent(c.Writer, gin.H{
		"type":    "CONNECTED",
		"message": "SSE connection established",
		"time":    time.Now(),
	})

	// Setup heartbeat
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Event loop
	for {
		select {
		case event := <-client.Events:
			if err := sendEvent(c.Writer, event); err != nil {
				return
			}
		case <-ticker.C:
			if err := sendEvent(c.Writer, gin.H{"type": "HEARTBEAT"}); err != nil {
				return
			}
		case <-clientGone:
			return
		}
	}
}

// Helper function to send SSE events
func sendEvent(w http.ResponseWriter, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		return err
	}

	w.(http.Flusher).Flush()
	return nil
}

// GetConnectedClients returns the number of connected clients
func GetConnectedClients(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connected_clients": Manager.GetClientCount(),
	})
}
