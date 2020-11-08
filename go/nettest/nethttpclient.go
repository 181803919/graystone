package main

import (
	"fflog"
	"net/http"
)

func main() {
	fflog.Open()
	defer fflog.Close()

	resp, err := http.Get("http://www.baidu.com")
	if err != nil{
		fflog.FFDebug("http.Get err = %v", err)
	}
	defer resp.Body.Close()

	fflog.FFDebug("Status = %v", resp.Status)
	fflog.FFDebug("StatusCode = %v", resp.StatusCode)
	fflog.FFDebug("Header = %v", resp.Header)

	buf := make([]byte, 4 * 1024)
	var tmp string
	for{
		n, err := resp.Body.Read(buf)
		if n == 0{
			fflog.FFDebug("read err = %v", err)
			break
		}
		tmp += string(buf[:n])
	}

	fflog.FFDebug("tmp=%v", tmp)
}
