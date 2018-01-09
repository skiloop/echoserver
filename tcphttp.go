package echoserver

import (
	"net"
	"github.com/google/logger"
	"github.com/skiloop/httplib"
)

func TcpHttpEchoServe(addr string) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Server error, could not start listener : %s", err.Error())
	}
	logger.Infof("Echo version %s, listening on %s", CommitSha, addr)
	n := int64(0)
	for {
		client, err := server.Accept()
		if err == nil {
			n += 1
			go func(n int64, c net.Conn) {
				logger.Infof("[%d] new connection", n)
				defer tcpClose(n, c)
				// write response header
				c.Write([]byte("HTTP/1.1 200 OK\r\nConnection: keep-alive\r\nContent-Type: text/html; charset=utf-8\r\n\r\n"))
				headers, err := httplib.HttpReadHeader(c)
				if err != nil {
					logger.Error(err)
					return
				}
				c.Write(headers)
				c.Write([]byte("\r\n"))
			}(n, client)
		} else {
			logger.Infof("Accept error : %s", err.Error())
			break
		}
	}

}
