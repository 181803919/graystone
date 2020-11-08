package main

import (
	"fflog"
	"net"
	"strings"
)

func HandleConn(conn net.Conn)  {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fflog.FFDebug("%s connect sucessful", addr)

	tm_buf := make([]byte, 2048)

	for{
		n, err := conn.Read(tm_buf)
		if err != nil{
			fflog.FFDebug( "err = %v", err)
		}

		fflog.FFDebug( "recv buf:%s", string(tm_buf[:n]))
		conn.Write([]byte(strings.ToUpper(string(tm_buf[:n]))))
	}
}

func main() {
	defer fflog.Close()
	fflog.Open()

	listen, err := net.Listen("tcp", ":18000")
	if err != nil{
		fflog.FFDebug("err = %v", err)
		return
	}
	defer listen.Close()

	for{
		conn, err := listen.Accept()
		if err != nil{
			fflog.FFDebug( "err=%v", err)
			return
		}

		go HandleConn(conn)
	}
}
