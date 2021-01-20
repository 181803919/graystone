package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fflog"
	"fmt"
	"log"
	"net/http"
)

func md5Sum(s string) string{
	m := md5.New()
	m.Write([]byte (s))
	return hex.EncodeToString(m.Sum(nil))
}

const uploadPath = "./tmp"

func myHttpDo (w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "stone program")
	w.Write([]byte("stone see you"))

	err := r.ParseForm()
	if err != nil{
		fflog.FFDebug("Data Error %v", err)
		w.Write([]byte("Data Error"))
	}

	formData := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&formData)

	dirName := fmt.Sprintf("%v", formData["dir"])
	fileMd5 := fmt.Sprintf("%v", formData["md5"])
	fileName := fmt.Sprintf("%v", formData["file"])
	pkgMd5 := fmt.Sprintf("%v", formData["sign"])

	md5CheckStr := md5Sum(dirName + "daklsgjlja2389173a21gasglkhk" + fileMd5)

	fflog.FFDebug("fileName:%s", fileName)
	fflog.FFDebug("dirName:%s", dirName)
	fflog.FFDebug("filemd5:%s", fileMd5)
	fflog.FFDebug("pkgmd5:%s", pkgMd5)
	fflog.FFDebug("mymd5:%s", md5CheckStr)

	w.Write([]byte("success"))
}

const maxUploadSize = 2 * 1024 * 2014 // 2 MB

func main() {
	http.HandleFunc("/upload", myHttpDo)

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Print("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading files.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}