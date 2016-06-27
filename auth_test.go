package appengine

import (
	"testing"

	"github.com/labstack/echo"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"
)

var testUsers = []*user.User{
	{Email: "kazu@example.com", AuthDomain: "example.com", Admin: false},
	{Email: "kazu@example.com", AuthDomain: "example.com", Admin: true},
}

func TestAuth(t *testing.T) {
	r := newTestRequest(t, "GET", "/", nil)
	c := newTestContext(r)
	h := Auth(Admin)(func(c echo.Context) error { return nil })
	err := h(c)
	if err != ErrUnauthroized {
		t.Errorf("err = %v; want %v", err, ErrUnauthroized)
	}

	aetest.Login(testUsers[0], r)
	c = newTestContext(r)
	err = h(c)
	if err != ErrForbidden {
		t.Errorf("err = %v; want %v", err, ErrForbidden)
	}

	aetest.Login(testUsers[1], r)
	c = newTestContext(r)
	err = h(c)
	if err != nil {
		t.Errorf("err = %v; want nil", err)
	}
}
