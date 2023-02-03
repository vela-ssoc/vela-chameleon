package proxy

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

/*
	chameleon.proxy{
		name = "xxxxx",
		bind = "tcp://127.0.0.1:3309",
		remote = "tcp://172.31.61.67:3389"
	}
*/
func newLuaProxyChameleon(L *lua.LState) int {
	cfg := newConfig(L)

	proc := L.NewVelaData(cfg.Name, proxyTypeOf)
	if proc.IsNil() {
		proc.Set(newProxyGo(cfg))
	} else {
		proc.Data.(*proxyGo).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(env vela.Environment, uv lua.UserKV) {
	xEnv = env
	uv.Set("proxy", lua.NewFunction(newLuaProxyChameleon))
}
