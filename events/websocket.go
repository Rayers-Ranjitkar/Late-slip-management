package events

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development - customize this for production
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	Conn    *websocket.Conn
	Send    chan []byte
	UserID  string
	IsAdmin bool
	Manager *ClientManager
	writeMu sync.Mutex // Add mutex for write synchronization
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, userID string, isAdmin bool, manager *ClientManager) *Client {
	return &Client{
		Conn:    conn,
		Send:    make(chan []byte, 256),
		UserID:  userID,
		IsAdmin: isAdmin,
		Manager: manager,
	}
}

// ReadPump handles incoming messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512) // Max message size
	c.Conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// We're not processing incoming messages in this implementation
		// but you could add message handling here if needed
	}
}

// WritePump handles sending messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel closed
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.writeMu.Lock() // Acquire the lock before writing
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.writeMu.Unlock() // Ensure the lock is released on error
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				c.writeMu.Unlock() // Ensure the lock is released on error
				return
			}
			c.writeMu.Unlock() // Release the lock after writing
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			// Send ping
			heartbeat := struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			}{
				Type: "HEARTBEAT",
				Data: nil,
			}

			heartbeatJSON, _ := json.Marshal(heartbeat)
			if err := c.Conn.WriteMessage(websocket.TextMessage, heartbeatJSON); err != nil {
				return
			}
		}
	}
}

// SendJSON sends a JSON message to the client
func (c *Client) SendJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	select {
	case c.Send <- jsonData:
		return nil
	default:
		return nil // Non-blocking send
	}
}
