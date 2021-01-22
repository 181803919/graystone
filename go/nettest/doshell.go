package main

import (
	"crypto/md5"
	"encoding/hex"
	"ffdaemon"
	"fflog"
	"net/http"
	"os/exec"
	"strings"
)

func doShell(w http.ResponseWriter, r *http.Request )  {
	r.ParseForm()
	cmdStr := r.FormValue("cmd")
	param1 := r.FormValue("param1")
	param2 := r.FormValue("param2")
	param3 := r.FormValue("param3")
	sign := r.FormValue("sign")

	m := md5.New()
	m.Write([]byte (cmdStr + "J2@1kjxl8Jlk##!"))
	if sign != hex.EncodeToString(m.Sum(nil)){
		fflog.FFError("Do Shell Fail PKG NOT VALID")
		w.Write([]byte ("Do Shell Fail PKG NOT VALID\n" ))
		return
	}

	cmd := exec.Command("/bin/bash", "-c",
		cmdStr + " " + param1 + " " + param2 + " " + param3)
	strRet, err := cmd.Output()
	if err != nil{
		fflog.FFError("Do Shell Fail" + err.Error())
		w.Write([]byte ("Do Shell Fail" + err.Error()))
		return
	}

	fflog.FFDebug("Do Shell Suc" + strings.Trim(string(strRet),"\n"))
	w.Write([]byte ("Do Shell Suc\n" + strings.Trim(string(strRet),"\n")))
}

func main(){
	ffdaemon.Daemon()
	fflog.Open()
	defer fflog.Close()

	fflog.FFDebug("Listen on Port 16667")
	http.HandleFunc("/doshell", doShell)

	err := http.ListenAndServe(":16667", nil)
	if err != nil{
		fflog.FFError("Listen 16667 fail" + err.Error())
		return
	}
}
