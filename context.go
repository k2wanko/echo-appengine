package appengine

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"google.golang.org/appengine"
)

func AppContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if r, ok := c.Request().(*standard.Request); ok {
				c.SetStdContext(appengine.WithContext(c.StdContext(), r.Request))
			}
			return next(c)
		}
	}
}
