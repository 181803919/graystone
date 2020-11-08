package main

import (
	"fflog"
	"net/http"
	"fmt"
)

func myHttpDo(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "stone program")
	w.Write([]byte("stone see you"))
	fflog.FFDebug("%v", r)
	fflog.FFDebug("r.Method = %v", r.Method)
	fflog.FFDebug("r.URL = %v", r.URL)
	fflog.FFDebug("r.Header = %v", r.Header)
	fflog.FFDebug("r.Body = %v", r.Body)
}

func main() {
	fflog.Open()
	defer fflog.Close()

	http.HandleFunc("/stonetest", myHttpDo)
	http.ListenAndServe(":18001", nil)
}
