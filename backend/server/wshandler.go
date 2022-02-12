package server

import (
	"errors"
	"net/http"

	chatwebsocket "github.com/alejoacosta74/chatserver/websocket"
	"github.com/labstack/echo"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: false,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(c echo.Context) error {
	// upgrade this connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if !ok {
		c.Logger().Error("could not find my ctx")
		return errors.New("could not find myctx")
	}
	cc.GetDebugLogger().Log("msg", "Got websocket request")

	pool := cc.GetPool()

	client := chatwebsocket.NewClient(conn, pool, cc.GetLogger(), cc.IsDebugEnabled())
	pool.Register <- client
	client.Read()

	return nil
}
