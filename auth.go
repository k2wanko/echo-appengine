package appengine

import (
	"net/http"

	"github.com/labstack/echo"
	"google.golang.org/appengine/user"
)

const (
	Optional = "optional"
	Required = "required"
	Admin    = "admin"
)

var (
	ErrUnauthroized = echo.NewHTTPError(http.StatusUnauthorized, "unauthroized")
	ErrForbidden    = echo.NewHTTPError(http.StatusForbidden, "forbidden")
)

func Auth(login string) echo.MiddlewareFunc {
	if login != Required && login != Admin {
		login = Optional
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := user.Current(c)
			//TODO: context.WithValue
			if login == Optional {
				return next(c)
			}
			if u == nil {
				return ErrUnauthroized
			}
			if login == Admin && !u.Admin {
				return ErrForbidden
			}
			return next(c)
		}
	}
}

func AuthErrorRedirect(err error, c echo.Context) {
	if err == ErrUnauthroized || err == ErrForbidden {
		url, err := user.LoginURL(c, c.Request().URL().Path())
		if err != nil {
			c.Error(err)
			return
		}
		c.Redirect(http.StatusFound, url)
	}
}
