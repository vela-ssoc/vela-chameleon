package stream

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

func (s *stream) pipeL(L *lua.LState) int {
	s.cfg.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (s *stream) startL(L *lua.LState) int {
	xEnv.Start(L, s).From(s.Code()).Do()
	return 0
}

func (s *stream) ignoreL(L *lua.LState) int {
	s.cur.ignore.CheckMany(L)
	return 0
}

func (s *stream) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return L.NewFunction(s.pipeL)
	case "ignore":
		return L.NewFunction(s.ignoreL)
	case "start":
		return L.NewFunction(s.startL)

	}
	return lua.LNil
}
