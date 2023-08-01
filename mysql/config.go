package mysql

import (
	"github.com/vela-ssoc/vela-chameleon/mysql/auth"
	"github.com/vela-ssoc/vela-chameleon/mysql/server"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/xreflect"
	"github.com/vela-ssoc/vela-kit/lua"
)

type config struct {
	Name     string    `lua:"name"     type:"string"`
	Bind     string    `lua:"bind"     type:"string"`
	Auth     auth.Auth `lua:"auth"     type:"object"`
	Database *EngineDB `lua:"database" type:"object"`

	CodeVM string
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := new(config)

	cfg.CodeVM = L.CodeVM()

	if e := xreflect.ToStruct(tab, cfg); e != nil {
		L.RaiseError("%v", e)
		return cfg
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return cfg
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := auxlib.Name(cfg.Name); e != nil {
		return e
	}

	return nil
}

func (cfg *config) toSerCfg() server.Config {
	return server.Config{
		Protocol: "tcp",
		Address:  cfg.Bind,
		Auth:     cfg.Auth,
		CodeVM:   func() string { return cfg.CodeVM },
	}
}
