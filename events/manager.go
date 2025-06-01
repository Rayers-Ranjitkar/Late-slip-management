package events

import (
	"sync"
)

// NotificationManager handles all SSE clients and broadcasting
type NotificationManager struct {
	clients map[string]*SSEClient
	mu      sync.RWMutex
}

// Instance of NotificationManager that will be used across the application
var Manager = NewManager()

// NewManager creates a new instance of NotificationManager
func NewManager() *NotificationManager {
	return &NotificationManager{
		clients: make(map[string]*SSEClient),
	}
}

// AddClient registers a new client for SSE updates
func (m *NotificationManager) AddClient(id string, isAdmin bool) *SSEClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	client := NewClient(id, isAdmin)
	m.clients[id] = client
	return client
}

// RemoveClient removes a client and closes their event channel
func (m *NotificationManager) RemoveClient(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[id]; exists {
		close(client.Events)
		delete(m.clients, id)
	}
}

// BroadcastEvent sends an event to all connected clients
func (m *NotificationManager) BroadcastEvent(event interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.clients {
		select {
		case client.Events <- event:
			// Event sent successfully
		default:
			// Channel is full, skip this client
		}
	}
}

// SendAdminNotification sends events to all admin clients
func (m *NotificationManager) SendAdminNotification(event interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.clients {
		if client.IsAdmin {
			select {
			case client.Events <- event:
				// Event sent successfully
			default:
				// Channel is full, skip this client
			}
		}
	}
}

// SendEventToClient sends an event to a specific client
func (m *NotificationManager) SendEventToClient(clientID string, event interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if client, exists := m.clients[clientID]; exists {
		select {
		case client.Events <- event:
			return true
		default:
			return false
		}
	}
	return false
}

// GetClientCount returns the number of connected clients
func (m *NotificationManager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}
