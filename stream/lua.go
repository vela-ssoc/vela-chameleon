package stream

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

/*
	chameleon.stream{
		name = "ssss",
		bind = "tcp://127.0.0.1:3390",
		remote = "tcp://172.31.61.67:3389",
	}
*/

func newLuaStreamChameleon(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewVelaData(cfg.name, streamTypeOf)
	if proc.IsNil() {
		proc.Set(newStream(cfg))
	} else {
		proc.Data.(*stream).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(env vela.Environment, uv lua.UserKV) {
	xEnv = env
	uv.Set("stream", lua.NewFunction(newLuaStreamChameleon))
}
