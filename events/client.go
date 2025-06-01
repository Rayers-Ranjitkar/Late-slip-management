package events

type SSEClient struct {
	ID      string
	Events  chan interface{}
	IsAdmin bool
}

func NewClient(id string, isAdmin bool) *SSEClient {
	return &SSEClient{
		ID:      id,
		Events:  make(chan interface{}, 10),
		IsAdmin: isAdmin,
	}
}
