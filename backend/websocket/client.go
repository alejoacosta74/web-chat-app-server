package websocket

import (
	// "github.com/alejoacosta74/chatserver/pkg/server"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     uuid.UUID       // a uniquely identifiably ID for a particular connection
	Conn   *websocket.Conn // a pointer to a websocket.Conn object
	Pool   *Pool           // a pointer to the Pool which this client will be part of
	logger log.Logger
	debug  bool
}

func NewClient(conn *websocket.Conn, pool *Pool, logger log.Logger, debug bool) *Client {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return &Client{
		ID:     id,
		Conn:   conn,
		Pool:   pool,
		logger: logger,
		debug:  debug,
	}
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			c.GetErrorLogger().Log("msgType", messageType, "body", string(p), "error", err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		c.GetDebugLogger().Log("msg", message.Body)
	}
}

func (c *Client) GetDebugLogger() log.Logger {
	if !c.debug {
		return log.NewNopLogger()
	}
	return log.With(level.Debug(c.logger), "component", "client", "caller", log.DefaultCaller)
}

func (c *Client) GetLogger() log.Logger {
	c.logger = log.WithPrefix(c.logger, "component", "client")
	return log.With(level.Info(c.logger))
}

func (c *Client) GetErrorLogger() log.Logger {
	return log.With(level.Error(c.logger), "component", "client", "caller", log.DefaultCaller)
}
