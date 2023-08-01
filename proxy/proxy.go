package proxy

import (
	"context"
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	risk "github.com/vela-ssoc/vela-risk"
	"net"
	"reflect"
	"time"
)

var proxyTypeOf = reflect.TypeOf((*proxyGo)(nil)).String()

type proxyGo struct {
	lua.SuperVelaData
	cfg *config
	cur config
	ln  *kind.Listener
}

func newProxyGo(cfg *config) *proxyGo {
	p := &proxyGo{cfg: cfg}
	p.V(lua.VTInit, proxyTypeOf)
	return p
}

func (p *proxyGo) Name() string {
	return p.cur.Name
}

func (p *proxyGo) Code() string {
	return p.cfg.co.CodeVM()
}

func (p *proxyGo) equal() bool {
	if p.cur.Bind.String() != p.cfg.Bind.String() {
		return false
	}

	if p.cur.Remote.String() != p.cfg.Remote.String() {
		return false
	}

	return true
}

func (p *proxyGo) Listen() error {

	if p.ln == nil {
		goto conn
	}

conn:
	ln, err := kind.Listen(xEnv, p.cfg.Bind)
	if err != nil {
		return err
	}
	p.ln = ln
	p.cur = *p.cfg

	return nil
}

func (p *proxyGo) Start() error {

	if e := p.Listen(); e != nil {
		return e
	}
	err := p.ln.OnAccept(p.accept)
	return err
}

//func (p *proxyGo) Reload() error {
//	return p.Listen()
//}

func (p *proxyGo) Close() error {
	e := p.ln.Close()
	return e
}

func (p *proxyGo) dail(conn net.Conn) (net.Conn, error) {

	host := p.cur.Remote.Hostname()
	port := p.cur.Remote.Port()

	if port == 0 {
		_, port = auxlib.ParseAddr(conn.LocalAddr())
	}

	if port == 0 {
		return nil, fmt.Errorf("invalid stream port")
	}

	d := net.Dialer{Timeout: 2 * time.Second}
	return d.Dial(p.cur.Remote.Scheme(), fmt.Sprintf("%s:%d", host, port))
}

func (p *proxyGo) pipe(ev *risk.Event) {
	if p.cur.log {
		ev.Log()
	}

	p.cur.pipe.Do(ev, p.cur.co, func(e error) {
		xEnv.Errorf("%s pipe call fail %v", p.Name(), e)
	})

	if p.cur.alert && ev.Alert {
		ev.Send()
	}

}

func (p *proxyGo) over(conn net.Conn) *risk.Event {
	ev := risk.HoneyPot()
	ev.Alert = false
	ev.Notice()
	ev.Subjectf("代理蜜罐请求结束")
	ev.Remote(conn)
	ev.From(p.cur.co)
	return ev
}

func (p *proxyGo) ignore(cnn net.Conn) bool {
	return p.cur.ignore.Match(kind.NewConn(cnn))
}

func (p *proxyGo) accept(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	if p.ignore(conn) {
		return nil
	}

	ev := risk.HoneyPot()
	ev.From(p.cur.co)
	ev.Remote(conn.RemoteAddr())

	dst, err := p.dail(conn)
	if err != nil {
		ev.Payloadf("%s 服务端口:%s 后端地址:%s 原因:%v",
			p.Name(), conn.LocalAddr().String(), p.cfg.Remote, err)
		p.pipe(ev)
		return err
	}
	defer dst.Close()

	ev.Payloadf("%s 服务端口:%s 后端地址:%s 链接成功", p.Name(), conn.LocalAddr().String(),
		dst.RemoteAddr().String())
	p.pipe(ev)

	newCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer func() {
			cancel()
		}()

		var toTn int64
		toTn, err = auxlib.Copy(newCtx, dst, conn)

		xEnv.Infof("%s 代理发送:%d 错误:%v", p.Name(), toTn, err)
	}()

	go func() {
		defer func() {
			cancel()
		}()
		var rev int64
		rev, err = auxlib.Copy(newCtx, conn, dst)
		xEnv.Infof("%s 接收远程:%d 错误:%v", p.Name(), rev, err)
	}()

	<-newCtx.Done()
	xEnv.Errorf("%s %s %s connection over.", p.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String())
	return err
}
