package echoserver

import (
	"net"
	"fmt"
	"sync"
	"github.com/google/logger"
)

type Handler interface {
	handle(conn net.Conn) (error)
}

type Listener struct {
	listener  net.Listener
	closeFlag bool
}
type Server struct {
	handler       Handler
	listener      net.Listener
	closeFlag     bool
	maxConcurrent int
	conn          net.Conn
	waitGroup     sync.WaitGroup
}

func New(handler Handler, maxConcurrent int) Server {
	return Server{handler: handler, maxConcurrent: maxConcurrent, closeFlag: false}
}
func (srv *Server) Start(bindAddr string, bindPort int64) (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	srv.listener = l
	srv.startWorkers()
	logger.Infof("server starts at %v", tcpAddr)
	cnt := 0
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			logger.Warning("accept error: %v", err)
			if srv.closeFlag {
				//close(srv.conn)
				return err
			}
			continue
		}
		logger.Infof("[%d] new connection", cnt)
		cnt++
		//srv.conn <- conn
		go srv.handler.handle(conn)
	}
	srv.waitGroup.Wait()
	return err
}

func (srv *Server) Close() {
	if !srv.closeFlag {
		srv.closeFlag = true
		srv.listener.Close()
	}
}
func (srv *Server) startWorkers() {
	for i := 0; i < srv.maxConcurrent; i++ {
		srv.waitGroup.Add(1)
		go srv.doWork(i)
	}
}

func (srv *Server) doWork(wid int) {
	logger.Infof("worker %d starts", wid)
	for ; !srv.closeFlag; {
		//conn, ok := <-srv.conn
		ok := false
		if !ok {
			logger.Infof("worker [%d] exit because of chan closed")
			break
		}
		if err := srv.handler.handle(srv.conn); err != nil {
			logger.Errorf("worker [%d] handle connection error: %v", wid, err)
		}
	}
	logger.Infof("worker [%d] done", wid)
	srv.waitGroup.Done()
}
