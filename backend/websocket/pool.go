package websocket

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

type Pool struct {
	Register   chan *Client     // register channel will send out 'New User Joined...' to all of the clients within this pool when a new client connects.
	Unregister chan *Client     // Will unregister a user and notify the pool when a client disconnects.
	Clients    map[*Client]bool // a map of clients to a boolean value. We can use the boolean value to dictate active/inactive
	Broadcast  chan Message     // when it is passed a message, will loop through all clients in the pool and send the message through the socket connection.
	logger     log.Logger
	debug      bool
}

func NewPool(logger log.Logger, debug bool) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		logger:     logger,
		debug:      debug,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			pool.GetLogger().Log("action", "register", "clientID", client.ID.String(), "poolsize", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			pool.GetLogger().Log("action", "unregister", "clientID", client.ID.String(), "poolsize", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
		case message := <-pool.Broadcast:
			pool.GetDebugLogger().Log("action", "broadcast", "message", message.Body)
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					pool.GetLogger().Log("broadcasting error", err)
					return
				}
			}
		}
	}
}

func (p *Pool) GetDebugLogger() log.Logger {
	if !p.debug {
		return log.NewNopLogger()
	}
	return log.With(level.Debug(p.logger), "caller", log.DefaultCaller)
}

func (p *Pool) GetLogger() log.Logger {
	return log.With(level.Info(p.logger), "component", "pool")
}
