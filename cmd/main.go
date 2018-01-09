package main

import (
	"github.com/skiloop/echoserver"
	"github.com/google/logger"
	"github.com/dmulholland/args"
	"os"
	"fmt"
)

func httpecho(p *args.ArgParser) {
	srv := echoserver.NewHttpServer(p.GetString("host"), int64(p.GetInt("port")))
	srv.Start()
}
func tcpecho(p *args.ArgParser) {
	addr := fmt.Sprintf("%s:%d", p.GetString("host"), int64(p.GetInt("port")))
	echoserver.TcpEchoServe(addr)
}
func tcphttpecho(p *args.ArgParser) {
	addr := fmt.Sprintf("%s:%d", p.GetString("host"), int64(p.GetInt("port")))
	echoserver.TcpHttpEchoServe(addr)
}

func main() {
	logger.Init("echoserver", false, false, os.Stdout)

	parser := args.NewParser()
	parser.Version = "0.0.1"
	parser.NewFlag("verbose v")
	cmdParser := parser.NewCmd("http", "http echo server", httpecho);
	cmdParser.NewString("host h", "127.0.0.1")
	cmdParser.NewInt("port p", 5000)

	cmdParser = parser.NewCmd("tcphttp", "tcp http echo server", tcphttpecho);
	cmdParser.NewString("host h", "127.0.0.1")
	cmdParser.NewInt("port p", 5000)
	cmdParser = parser.NewCmd("tcp", "tcp echo server", tcpecho);
	cmdParser.NewString("host h", "127.0.0.1")
	cmdParser.NewInt("port p", 5000)
	parser.Parse()
}
