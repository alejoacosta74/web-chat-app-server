package server

import (
	"io"

	"github.com/alejoacosta74/chatserver/websocket"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/labstack/echo"
)

type myCtx struct {
	echo.Context
	logWriter io.Writer
	logger    log.Logger
	debug     bool
	pool      *websocket.Pool
}

func (c *myCtx) SetLogWriter(logWriter io.Writer) {
	c.logWriter = logWriter
}

func (c *myCtx) GetLogWriter() io.Writer {
	return c.logWriter
}

func (c *myCtx) SetLogger(l log.Logger) {
	c.logger = log.WithPrefix(l, "component", "context")
}

func (c *myCtx) GetLogger() log.Logger {
	return c.logger
	// c.logger = log.WithPrefix(c.logger, "component", "context")
	// return log.With(level.Info(c.logger))
}

func (c *myCtx) GetDebugLogger() log.Logger {
	if !c.debug {
		return log.NewNopLogger()
	}
	return log.With(level.Debug(c.logger), "component", "myCTX", "caller", log.DefaultCaller)
}

func (c *myCtx) GetErrorLogger() log.Logger {
	return log.With(level.Error(c.logger), "component", "myCTX", "caller", log.DefaultCaller)
}

func (c *myCtx) GetPool() *websocket.Pool {
	return c.pool
}

func (c *myCtx) IsDebugEnabled() bool {
	return c.debug
}
