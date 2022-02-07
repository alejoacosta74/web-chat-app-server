package server

import (
	"errors"
	"net/http"

	// "github.com/alejoacosta74/chatserver/server/middlewares"
	// "github.com/gofrs/uuid"
	"github.com/labstack/echo"
)

func httpHandler(c echo.Context) error {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if !ok {
		c.Logger().Error("could not find my ctx")
		return errors.New("could not find myctx")
	}
	cc.GetDebugLogger().Log("msg", "Got HTTP request")
	return c.String(http.StatusOK, "Real time chat server")

}

// HealthCheck - Health Check Handler
func healthCheck(c echo.Context) error {
	// if requestID, ok := c.Get(middlewares.RequestIDContextKey).(uuid.UUID); ok {
	// 	c.Logger().Infof("RequestID: %s", requestID)
	// }
	resp := struct{ Message string }{
		Message: "Everything is good!",
	}
	return c.JSON(http.StatusOK, resp)
}
