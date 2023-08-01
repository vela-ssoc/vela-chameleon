package ssh

import (
	"context"
	"errors"
	"github.com/vela-ssoc/vela-kit/lua"
	risk "github.com/vela-ssoc/vela-risk"
)

type sshGo struct {
	lua.SuperVelaData
	cfg  *config
	auth *auth

	ctx    context.Context
	cancel context.CancelFunc

	serv *Server
}

func newSSH(cfg *config) *sshGo {
	s := &sshGo{cfg: cfg}
	s.V(lua.VTInit, sshTypeOf)
	s.auth = &auth{data: make(map[string]string)}
	return s
}

func (s *sshGo) Name() string {
	return s.cfg.name
}

func (s *sshGo) event(ctx Context, pass string, err error) {
	ev := risk.HoneyPot()
	ev.Subject = "ssh蜜罐认证失败"
	ev.From(s.cfg.code)
	ev.Remote(ctx.RemoteAddr())
	ev.Payloadf("user:%s pass:%s", ctx.User(), pass)
	ev.Send()
}

var (
	invalidU = errors.New("not found user")
	invalidP = errors.New("invalid pass")
)

func (s *sshGo) doAuth(ctx Context, pass string) bool {
	var err error
	defer s.event(ctx, pass, err)

	name := ctx.User()
	v, ok := s.auth.Get(name)
	if !ok {
		err = errors.New("not found user")
		goto ERR
	}

	if v != pass {
		err = errors.New("invalid pass")
		goto ERR
	}

ERR:
	return false
}

func (s *sshGo) handler(sess Session) {
}

func (s *sshGo) Start() error {
	s.serv = s.cfg.toSSH(s.handler, s.doAuth)
	s.serv.Version = s.cfg.version
	s.serv.CodeVM = func() string {
		return s.cfg.code
	}

	var err error
	xEnv.Spawn(100, func() {
		err = s.serv.ListenAndServe()
	})

	if err != nil {
		return err
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	xEnv.Errorf("%s %s start succeed", s.Name(), s.Type())
	return err
}

func (s *sshGo) Close() error {
	s.cancel()
	return s.serv.Close()
}
