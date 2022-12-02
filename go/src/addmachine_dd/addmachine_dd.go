package main

import (
	"ffdaemon"
	"fflog"
	"net/http"
)

func modifyMachine(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method == "POST" {
		fflog.FFDebug("after decrypt msg: POST")
	}
}

func main() {
	//test git
	ffdaemon.Daemon()
	fflog.Open()
	defer fflog.Close()

	fflog.FFDebug("Listen on Port 16669")
	http.HandleFunc("/modifyMachine", modifyMachine)

	err := http.ListenAndServe(":16669", nil)
	if err != nil {
		fflog.FFError("Listen 16669 fail" + err.Error())
		return
	}
}
