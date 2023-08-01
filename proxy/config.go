package proxy

import (
	"errors"
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

type config struct {
	Name   string
	Bind   auxlib.URL
	Remote auxlib.URL
	alert  bool
	log    bool
	ignore *cond.Ignore //ignore by s_addr , s_port , d_addr , d_port , source , destination
	filter *cond.Combine
	pipe   *pipe.Chains
	co     *lua.LState
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{
		ignore: cond.NewIgnore(),
		filter: cond.NewCombine(),
		pipe:   pipe.New(pipe.Env(xEnv)),
	}

	tab.Range(func(k string, v lua.LValue) {
		switch k {
		case "name":
			cfg.Name = auxlib.CheckProcName(v, L)
		case "alert":
			cfg.alert = lua.IsTrue(v)
		case "log":
			cfg.log = lua.IsTrue(v)
		case "bind":
			cfg.Bind = auxlib.CheckURL(v, L)
		case "remote":
			cfg.Remote = auxlib.CheckURL(v, L)
		}
	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}
	cfg.co = xEnv.Clone(L)
	return cfg
}

func (cfg *config) verify() error {
	if e := auxlib.Name(cfg.Name); e != nil {
		return e
	}

	if cfg.Bind.IsNil() {
		return fmt.Errorf("not found bind url")
	}

	if cfg.Remote.IsNil() {
		return fmt.Errorf("not found remote url")
	}

	switch cfg.Bind.Scheme() {
	case "tcp", "udp":
		//todo

	default:
		return errors.New("invalid bind protocol")
	}

	switch cfg.Remote.Scheme() {
	case "tcp", "udp":
		//todo

	default:
		return errors.New("invalid bind protocol")
	}

	return nil
}
