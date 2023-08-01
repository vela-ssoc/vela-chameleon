package proxy

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

func (p *proxyGo) pipeL(L *lua.LState) int {
	p.cfg.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (p *proxyGo) startL(L *lua.LState) int {
	xEnv.Start(L, p).From(p.CodeVM()).Do()
	return 0
}

func (p *proxyGo) ignoreL(L *lua.LState) int {
	p.cur.ignore.CheckMany(L)
	return 0
}

func (p *proxyGo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return L.NewFunction(p.pipeL)
	case "start":
		return L.NewFunction(p.startL)
	case "ignore":
		return L.NewFunction(p.ignoreL)
	}

	return lua.LNil
}
