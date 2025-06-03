package events

import (
	"time"

	"github.com/gin-gonic/gin"
)

var clientManager = NewClientManager()

// SSEHandler handles SSE connections for both admin and student clients
func SSEHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	isAdmin := role == "admin"
	messageChan := clientManager.AddClient(userID, isAdmin)
	defer clientManager.RemoveClient(userID, messageChan, isAdmin)

	// Set SSE headers
	SetHeaders(c)

	// Send initial connection message
	SendEventWithType(c, "CONNECTED", gin.H{
		"message": "SSE connection established",
		"time":    time.Now(),
	})

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	clientGone := c.Writer.CloseNotify()

	for {
		select {
		case msg := <-messageChan:
			if err := SendMessage(c, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := SendEventWithType(c, "HEARTBEAT", nil); err != nil {
				return
			}
		case <-clientGone:
			return
		}
	}
}
