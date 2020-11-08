package main

import (
	"fflog"
	"net"
	"os"
)

func main() {
	defer fflog.Close()
	fflog.Open()

	conn, err := net.Dial("tcp", "127.0.0.1:18000")
	if err != nil{
		fflog.FFError("net.Dial %v", err)
		return
	}
	defer conn.Close()

	go func() {
		str := make([] byte, 1024)
		for{
			n, err := os.Stdin.Read(str)
			if err != nil {
				fflog.FFError("conn.Read err = %v", err)
			}
			conn.Write(str[:n])
		}
	}()

	buf := make([]byte, 1024)
	for{
		n, err := conn.Read(buf)
		if err != nil{
			fflog.FFError("conn.Read err:%v", err)
			return
		}
		fflog.FFDebug("Recv:%s", string(buf[:n]))
	}
}
