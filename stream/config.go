package stream

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

type config struct {
	name   string
	bind   auxlib2.URL
	remote auxlib2.URL

	alert  bool
	log    bool
	ignore *cond.Ignore
	filter *cond.Combine
	pipe   *pipe.Chains
	co     *lua.LState
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{
		pipe: pipe.New(),
	}

	tab.Range(func(key string, lv lua.LValue) {
		switch key {
		case "name":
			cfg.name = auxlib2.CheckProcName(lv, L)
		case "bind":
			cfg.bind = auxlib2.CheckURL(lv, L)
		case "remote":
			cfg.remote = auxlib2.CheckURL(lv, L)
		case "alert":
			cfg.alert = lua.IsTrue(lv)

		default:
			//todo
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
	if err := auxlib2.Name(cfg.name); err != nil {
		return err
	}

	switch cfg.bind.Scheme() {
	case "tcp", "udp", "unix":
		return nil
	default:
		return fmt.Errorf("%s invalid bind url", cfg.name)
	}

	switch cfg.remote.Scheme() {
	case "tcp", "udp", "unix", "http":
		return nil
	default:
		return fmt.Errorf("%s invalid remote url", cfg.name)
	}

	return nil
}
