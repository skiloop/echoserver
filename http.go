package echoserver

import (
	"net/http"
	"bytes"
	"fmt"
)

type HttpServer struct {
	addr string
	port int64
}

func (srv *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBufferString(fmt.Sprintf("%s %s %s\n", r.Method, r.URL.Path, r.Proto))

	for k, v := range r.Header {
		s := bytes.NewBufferString("")
		s.WriteString(v[0])
		for vv := range v[1:] {
			s.WriteString("," + v[vv])
		}
		buf.WriteString(fmt.Sprintf("%v: %v\n", k, s))
	}
	buf.ReadFrom(r.Body)
	w.Write(buf.Bytes())
}

func (srv *HttpServer) Start() {
	http.ListenAndServe(fmt.Sprintf("%v:%d", srv.addr, srv.port), srv)
}

func NewHttpServer(addr string, port int64) (srv *HttpServer) {
	srv = &HttpServer{addr: addr, port: port}
	return srv
}
