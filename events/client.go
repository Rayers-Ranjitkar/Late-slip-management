package events

import "sync"

type SSEClient struct {
	ID      string
	Events  chan string
	IsAdmin bool
}

type ClientManager struct {
	adminClients   map[chan string]bool
	studentClients map[string]chan string
	mu             sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		adminClients:   make(map[chan string]bool),
		studentClients: make(map[string]chan string),
	}
}

func (cm *ClientManager) AddClient(userID string, isAdmin bool) chan string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	messageChan := make(chan string, 10)
	if isAdmin {
		cm.adminClients[messageChan] = true
	} else {
		cm.studentClients[userID] = messageChan
	}
	return messageChan
}

func (cm *ClientManager) RemoveClient(userID string, messageChan chan string, isAdmin bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if isAdmin {
		delete(cm.adminClients, messageChan)
	} else {
		delete(cm.studentClients, userID)
	}
	close(messageChan)
}
