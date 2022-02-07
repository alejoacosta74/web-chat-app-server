package server

import (
	"fmt"
	"io"
	"sync"

	// chatmiddleware "github.com/alejoacosta74/chatserver/server/middlewares"
	"github.com/alejoacosta74/chatserver/websocket"
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	address   string
	logWriter io.Writer
	logger    log.Logger
	debug     bool
	mutex     *sync.Mutex
	echo      *echo.Echo
}

func New(address string, opts ...Option) (*Server, error) {
	s := &Server{
		logger:  log.NewNopLogger(),
		echo:    echo.New(),
		address: address,
	}

	var err error
	for _, opt := range opts {
		if err = opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Server) Start() error {
	logWriter := s.logWriter
	e := s.echo
	if s.mutex == nil {
		e.GET("/", httpHandler)
	} else {
		e.GET("/", func(c echo.Context) error {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			return httpHandler(c)
		})
	}
	e.GET("/ws", wsHandler)
	e.GET("/health-check", healthCheck)

	e.Use(middleware.BodyDump(func(c echo.Context, req []byte, res []byte) {
		myctx := c.Get("myctx")
		cc, ok := myctx.(*myCtx)
		if !ok {
			return
		}

		if s.debug {
			cc.GetDebugLogger().Log("msg", "HTTP")
			fmt.Fprintf(logWriter, "=> HTTP request\n%s\n", req)
			fmt.Fprintf(logWriter, "<= HTTP response\n%s\n", res)
		}
	}))

	pool := websocket.NewPool(s.logger, s.debug)
	go pool.Start()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &myCtx{
				Context:   c,
				logger:    s.logger,
				pool:      pool,
				logWriter: logWriter,
				debug:     s.debug,
			}

			c.Set("myctx", cc)
			return next(c)
		}
	})

	return e.Start(s.address)
}

type Option func(*Server) error

func SetLogWriter(logWriter io.Writer) Option {
	return func(s *Server) error {
		s.logWriter = logWriter
		return nil
	}
}

func SetLogger(l log.Logger) Option {
	return func(s *Server) error {
		s.logger = l
		return nil
	}
}

func SetDebug(debug bool) Option {
	return func(s *Server) error {
		s.debug = debug
		return nil
	}
}

func SetSingleThreaded(singleThreaded bool) Option {
	return func(s *Server) error {
		if singleThreaded {
			s.mutex = &sync.Mutex{}
		} else {
			s.mutex = nil
		}
		return nil
	}
}
