package chameleon

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-chameleon/mysql"
	"github.com/vela-ssoc/vela-chameleon/proxy"
	"github.com/vela-ssoc/vela-chameleon/ssh"
	"github.com/vela-ssoc/vela-chameleon/stream"
	"github.com/vela-ssoc/vela-kit/lua"
)

func WithEnv(env vela.Environment) {
	uv := lua.NewUserKV()
	proxy.Inject(env, uv)
	stream.Inject(env, uv)
	mysql.Inject(env, uv)
	ssh.Inject(env, uv)
	env.Global("chameleon", uv)
}
