package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"ffdaemon"
	"fflog"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

/*
test cmd:
curl -X POST "http://192.168.56.101:39100/svrlistupload" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "file=@/home/temp/stontestsvrlist.txt" -F 'data={"dir":"/home/client/","md5":"3jlsjidffd","file":"svrlist.txt","sign":"9jjkk332222"}'
*/

func main(){
	ffdaemon.Daemon()
	fflog.Open()
	//fflog.OpenEx("/home/golog/filesvrd", fflog.LOG_DEBUG)
	defer fflog.Close()

	fflog.FFDebug("Listen On 16666!")
	http.HandleFunc("/svrlistupload", upLoad)

	err := http.ListenAndServe("0.0.0.0:16666", nil)
	if err != nil{
		fflog.FFError("Listen Error:" + err.Error())
	}
}

type UploadExtendInfo struct{
	dirName string
	fileMd5 string
	fileName string
	pkgMd5 string
}

func myMd5Sum(srcStr string) string{
	m := md5.New()
	m.Write([]byte (srcStr))
	return hex.EncodeToString(m.Sum(nil))
}

func myMd5FileSum(file *multipart.File) string{
	m := md5.New()
	io.Copy(m, *file)
	return hex.EncodeToString(m.Sum(nil))
}

func initUpLoadInfo(uploadExtend *UploadExtendInfo, r *map[string][]string){
	tmpInfo := make(map[string]string)
	if r != nil{
		for formName, values := range *r{
			for _, value := range values{
				tmpInfo[formName] = value
			}
		}
	}

	uploadExtend.dirName =  tmpInfo["dir"]
	uploadExtend.fileMd5 = tmpInfo["md5"]
	uploadExtend.fileName = tmpInfo["file"]
	uploadExtend.pkgMd5 =  tmpInfo["sign"]
}

func initUpLoadInfoJson(uploadExtend *UploadExtendInfo, r *map[string][]string){
	if r != nil{
		for _, values := range *r{
			for _, value := range values{
				m := make(map[string]interface{})
				json.NewDecoder(strings.NewReader(value)).Decode(&m)
				uploadExtend.dirName = fmt.Sprintf("%v", m["dir"])
				uploadExtend.fileMd5 = fmt.Sprintf("%v", m["md5"])
				uploadExtend.fileName = fmt.Sprintf("%v", m["file"])
				uploadExtend.pkgMd5 = fmt.Sprintf("%v", m["sign"])
				fflog.FFDebug("fileName:%s", uploadExtend.fileName)
				fflog.FFDebug("dirName:%s", uploadExtend.dirName)
				fflog.FFDebug("filemd5:%s", uploadExtend.fileMd5)
				fflog.FFDebug("pkgmd5:%s", uploadExtend.pkgMd5)
			}
		}
	}

}

func upLoad(w http.ResponseWriter, r *http.Request)  {
	contentType := r.Header.Get("content-type")
	acceptType := r.Header.Get("accept")
	contentLen := r.ContentLength

	if !strings.Contains(contentType, "multipart/form-data"){
		fflog.FFError("content-type must be multipart/form-data")
		w.Write([]byte("content-type must be multipart/form-data\n"))
		return
	}

	if contentLen > 4*1024*1024 {
		fflog.FFError("File Size Max Limit 4M")
		w.Write([]byte("File Size Max Limit 4M\n"))
		return
	}

	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil{
		fflog.FFError("ParseMultipartForm Error" + err.Error())
		w.Write([]byte("ParseMultipartForm Error" + err.Error() + "\n"))
		return
	}

	var uploadExtend UploadExtendInfo

	if !strings.Contains(acceptType, "application/json"){
		initUpLoadInfo(&uploadExtend, &r.MultipartForm.Value)
	}else{
		initUpLoadInfoJson(&uploadExtend, &r.MultipartForm.Value)
	}

	myMd5Str := myMd5Sum(uploadExtend.dirName +
		"daklsgjlja2389173a21gasglkhk" +
		uploadExtend.fileMd5)
	if myMd5Str != uploadExtend.pkgMd5{
		fflog.FFError("Key Error")
		w.Write([]byte("Key Error\n"))
		return
	}

	if len(r.MultipartForm.File) <= 0{
		fflog.FFError("not have any file")
		w.Write([]byte("not have any file\n"))
		return
	}

	for name, files := range r.MultipartForm.File{
		fflog.FFDebug("upload file " + name)
		if len(files) != 1{
			fflog.FFError("too many file")
			w.Write([]byte("too many file\n"))
			return
		}
		if name == ""{
			fflog.FFError("is not FileData")
			w.Write([]byte("too many file\n"))
			return
		}
		for _, f := range files{
			handle, err := f.Open()
			if err != nil{

				fflog.FFError("unknown error, fileName:" + f.Filename +
					", size:" + strconv.FormatInt(f.Size, 10) + ", err:" +
					err.Error())
				w.Write([]byte("unknown error, fileName:" + f.Filename +
					", size:" + strconv.FormatInt(f.Size, 10) + ", err:" +
					err.Error() + "\n"))
				return
			}

			//fflog.FFDebug("mymd5:" + myMd5FileSum(&handle))
			if uploadExtend.fileMd5 != myMd5FileSum(&handle){
				fflog.FFError("File CHECK ERROR")
				w.Write([]byte("File CHECK ERROR\n"))
				return
			}
			handle.Seek(0, 0)

			os.MkdirAll(uploadExtend.dirName, 0777)
			wholePath := uploadExtend.dirName + path.Base(f.Filename)
			dst, _ := os.Create(wholePath)
			io.Copy(dst, handle)
			dst.Close()

			fflog.FFDebug("upload suc, fileName:" + f.Filename +
				", size:" + strconv.FormatInt(f.Size, 10))
			w.Write([]byte("upload suc, fileName:" + f.Filename +
				", size:" + strconv.FormatInt(f.Size, 10) + "\n"))
			return
		}
	}
}


