package stream

import (
	"context"
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	risk "github.com/vela-ssoc/vela-risk"
	"net"
	"reflect"
	"time"
)

var (
	streamTypeOf = reflect.TypeOf((*stream)(nil)).String()
)

type stream struct {
	lua.SuperVelaData
	cfg *config
	cur config //保存当前启动 为了下次快速启动
	ln  *kind.Listener
}

func newStream(cfg *config) *stream {
	obj := &stream{cfg: cfg}
	obj.V(lua.VTInit, streamTypeOf)
	return obj
}

func (s *stream) http() (vela.HTTPStream, error) {
	/*
		return xEnv.Stream("http", map[string]interface{}{
			"address": s.cur.remote.String(),
		})

	*/
	return nil, nil
}

func (s *stream) socket(conn net.Conn) (vela.HTTPStream, error) {
	/*

		switch s.cur.remote.Scheme() {
		case "http", "https":
			return s.http()

		default:
			host := s.cur.remote.Hostname()
			port := s.cur.remote.Port()
			if port == 0 {
				_, port = auxlib2.ParseAddr(conn.LocalAddr())
			}

			if port == 0 {
				return nil, fmt.Errorf("invalid stream port")
			}
			return xEnv.Stream("tunnel", map[string]interface{}{
				"network": s.cur.remote.Scheme(),
				"address": fmt.Sprintf("%s:%d", host, port),
			})

		}

	*/
	return nil, nil
}

func (s *stream) pipe(ev *risk.Event) {
	if s.cur.log {
		ev.Log()
	}

	s.cur.pipe.Do(ev, s.cur.co, func(e error) {
		xEnv.Errorf("%s stream pipe fail %v", s.Name(), e)
	})

	if s.cur.alert && ev.Alert {
		ev.Send()
	}
}

func (s *stream) Code() string {
	return s.cfg.co.CodeVM()
}

func (s *stream) ignore(conn net.Conn) bool {
	if s.cur.ignore == nil {
		return false
	}

	return s.cur.ignore.Match(kind.NewConn(conn))
}

func (s *stream) accept(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	if s.ignore(conn) {
		return nil
	}

	//toT nt
	bind := s.cur.bind.String()
	remote := s.cur.remote.String()

	ev := risk.HoneyPot()
	ev.From(s.CodeVM())
	ev.Remote(conn.RemoteAddr())
	ev.Local(conn.LocalAddr())
	ev.Subjectf("流式代理蜜罐命中")
	ev.Payloadf("后端%s", bind)
	s.pipe(ev)

	var toTn int64

	//接收的数据
	var rev int64

	//报错
	var err error

	//数据通道
	var socket vela.HTTPStream

	socket, err = s.socket(conn)
	if err != nil {
		return err
	}
	defer socket.Close()

	newCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer func() {
			cancel()
		}()

		toTn, err = auxlib2.Copy(newCtx, socket, conn)
		xEnv.Infof("stream %s proxy send %v data:%d", s.Name(), remote, toTn)
	}()

	go func() {
		defer func() {
			cancel()
		}()
		rev, err = auxlib2.Copy(ctx, conn, socket)
		xEnv.Infof("stream %s proxy recv  %s data:%d", s.Name(), remote, rev)
	}()

	<-newCtx.Done()
	xEnv.Errorf("%s %s-%s close", s.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String())
	return err
}

func (s *stream) equal() bool {
	if s.cfg.remote.String() != s.cur.remote.String() {
		return false
	}

	if s.cfg.bind.String() != s.cur.bind.String() {
		return false
	}

	return true

}

func (s *stream) Listen() error {

	if s.ln == nil {
		goto conn
	}

conn:
	ln, err := kind.Listen(xEnv, s.cfg.bind)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}

func (s *stream) start() (err error) {
	s.cur = *s.cfg
	err = s.ln.OnAccept(s.accept)
	<-time.After(time.Millisecond * 50)
	return err
}

func (s *stream) Start() error {

	if e := s.Listen(); e != nil {
		return e
	}

	return s.start()
}

//func (st *stream) Reload() (err error) {
//	if e := st.Listen(); e != nil {
//		return e
//	}
//
//	st.cur = *st.cfg
//	return nil
//}

func (s *stream) Close() error {
	return s.ln.Close()
}

func (s *stream) Name() string {
	return s.cur.name
}

func (s *stream) Type() string {
	return streamTypeOf
}
