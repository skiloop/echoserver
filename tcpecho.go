package echoserver

import (
	"net"
	"github.com/google/logger"
	"io"
)

var CommitSha string

func tcpClose(n int64, c net.Conn) {
	logger.Infof("[%d] connection close", n)
	c.Close()
}

func TcpEchoServe(addr string) {
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
				io.Copy(c, c)
			}(n, client)
		} else {
			logger.Infof("Accept error : %s", err.Error())
			break
		}
	}

}
