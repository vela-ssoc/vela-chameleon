package mysql

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-chameleon/mysql/auth"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

func newLuaMySQL(L *lua.LState) int {
	cfg := newConfig(L)

	proc := L.NewVelaData(cfg.Name, TGoMySQL)
	if proc.IsNil() {
		proc.Set(newGoMysql(cfg))
	} else {
		proc.Data.(*GoMysql).cfg = cfg
	}
	L.Push(proc)
	return 1
}

func newLuaAuth(L *lua.LState) int {
	name := L.CheckString(1)
	pass := L.CheckString(2)

	a := new(Audit)
	a.CodeVM = L.CodeVM

	obj := auth.NewAudit(auth.NewNativeSingle(name, pass, auth.AllPermissions), a)
	L.Push(L.NewAnyData(obj))
	return 1
}

func newLuaTable(L *lua.LState) int {
	tab := newTable(L)
	L.Push(L.NewAnyData(tab))
	return 1
}

func newLuaEngineDB(L *lua.LState) int {
	name := L.CheckString(1)
	db := newEngineDB(name)
	L.Push(L.NewAnyData(db))
	return 1
}

func Inject(env vela.Environment, uv lua.UserKV) {
	xEnv = env

	m := lua.NewUserKV()
	m.Set("new", lua.NewFunction(newLuaMySQL))
	m.Set("auth", lua.NewFunction(newLuaAuth))
	m.Set("db", lua.NewFunction(newLuaEngineDB))
	m.Set("table", lua.NewFunction(newLuaTable))
	uv.Set("mysql", m)
}
