package main

type Hub struct {
	Clients   map[*Client]bool
	Add       chan *Client
	Remove    chan *Client
	Broadcast chan []byte
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Add:
			h.Clients[client] = true
		case client := <-h.Remove:
			delete(h.Clients, client)
			client.Cleanup()
		case message := <-h.Broadcast:
			for c := range h.Clients {
				c.Send <- message
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*Client]bool),
		Add:       make(chan *Client),
		Remove:    make(chan *Client),
		Broadcast: make(chan []byte),
	}
}
