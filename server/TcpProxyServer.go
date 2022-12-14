package server

import (
	"context"
	"godeliver/conf"
	"godeliver/driver"
	"godeliver/misc"
	"log"
	"net"
	"time"
)

type TcpProxyServer struct {
	*driver.ProxyServer
}

func CreateTcpProxyServer() *TcpProxyServer {
	port := conf.GetIniValue("srv", "port")

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	log.Printf("start proxy on: %s", port)
	srv := &TcpProxyServer{}
	srv.ProxyServer = &driver.ProxyServer{}
	srv.Listener = listen
	return srv
}

func (this *TcpProxyServer) Run() {
	for {
		s, err := this.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}

func (this *TcpProxyServer) handleClt(clt net.Conn) {
	defer func() {
		clt.Close()
		log.Printf("clt disconnect %s\n", clt.RemoteAddr().String())
	}()

	clt.SetDeadline(time.Now().Add(30*time.Second))

	log.Printf("new clt: %s", clt.RemoteAddr().String())

	proxy, err := net.Dial("tcp", conf.GetIniValue("proxy", "ip")+":"+conf.GetIniValue("proxy", "port"))
	if err != nil {
		log.Println("err :", err)
		return
	}
	defer func() {
		proxy.Close()
		log.Printf("proxy disconnect %s\n", proxy.RemoteAddr().String())
	}() // 关闭连接


	ctx, cancel := context.WithCancel(context.Background())
	go this.HandlePtoC(ctx, clt, proxy)
	this.handleCtoP(ctx, clt, proxy)
	cancel()
}

func (this *TcpProxyServer) handleCtoP(ctx context.Context, clt net.Conn, proxy net.Conn) {
	for {
		clt.SetDeadline(time.Now().Add(30*time.Second))
		_, msg, err := misc.Recv(clt)
		if err != nil {
			log.Println(err)
			break
		}

		msg = this.Crypt.Decode(msg)
		_, err = proxy.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
		//log.Printf("send proxy[%s] %s\n", proxy.RemoteAddr().String(), msg)
	}
}