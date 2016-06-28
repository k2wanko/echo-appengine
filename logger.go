package appengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/labstack/echo"
	"github.com/labstack/echo/log"
	glog "github.com/labstack/gommon/log"
	"golang.org/x/net/context"
	aelog "google.golang.org/appengine/log"
)

type logger struct {
	ctx context.Context
	lv  glog.Lvl
	b   *bytes.Buffer
}

func (l *logger) json(j glog.JSON) (res string) {
	json.NewEncoder(l.b).Encode(j)
	res = l.b.String()
	l.b.Reset()
	return
}

func (l *logger) SetOutput(w io.Writer) {
	panic("echo.Logger: unsupport SetOutput")
}

func (l *logger) SetLevel(level glog.Lvl) {
	l.lv = level
}

func (l *logger) Print(args ...interface{}) {
	l.Printf("%s", fmt.Sprintln(args...))
}

func (l *logger) Printf(format string, args ...interface{}) {
	aelog.Infof(l.ctx, format, args...)
}

func (l *logger) Printj(j glog.JSON) {
	l.Printf("%s", l.json(j))
}

func (l *logger) Info(args ...interface{}) {
	l.Infof("%s", fmt.Sprintln(args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	if glog.INFO < l.lv {
		return
	}
	aelog.Infof(l.ctx, format, args...)
}

func (l *logger) Infoj(j glog.JSON) {
	l.Infof("%s", l.json(j))
}

func (l *logger) Debug(args ...interface{}) {
	l.Debugf("%s", fmt.Sprintln(args...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if glog.DEBUG < l.lv {
		return
	}
	aelog.Debugf(l.ctx, format, args...)
}

func (l *logger) Debugj(j glog.JSON) {
	l.Debugf("%s", l.json(j))
}

func (l *logger) Warn(args ...interface{}) {
	l.Warnf("%s", fmt.Sprintln(args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if glog.WARN < l.lv {
		return
	}
	aelog.Warningf(l.ctx, format, args...)
}

func (l *logger) Warnj(j glog.JSON) {
	l.Warnf("%s", l.json(j))
}

func (l *logger) Error(args ...interface{}) {
	l.Errorf("%s", fmt.Sprintln(args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if glog.ERROR < l.lv {
		return
	}
	aelog.Errorf(l.ctx, format, args...)
}

func (l *logger) Errorj(j glog.JSON) {
	l.Errorf("%s", l.json(j))
}

func (l *logger) Fatal(args ...interface{}) {
	l.Fatalf("%s", fmt.Sprintln(args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if glog.FATAL < l.lv {
		return
	}
	aelog.Criticalf(l.ctx, format, args...)
}

func (l *logger) Fatalj(j glog.JSON) {
	l.Fatalf("%s", l.json(j))
}

func Logger(ctx context.Context) log.Logger {
	return &logger{
		ctx: ctx,
		b:   new(bytes.Buffer),
		lv:  glog.INFO,
	}
}

func AppLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		var l log.Logger
		return func(c echo.Context) error {
			if l == nil {
				l = Logger(c)
				c.Echo().SetLogger(l)
			}
			return next(c)
		}
	}
}
