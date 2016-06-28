package appengine

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type InitOption func(*options)

type options struct {
	path string
}

func Init(e *echo.Echo, opts ...InitOption) {
	op := &options{
		path: "/",
	}

	for _, o := range opts {
		o(op)
	}

	s := standard.New("")
	s.SetHandler(e)
	http.Handle(op.path, s)
}
