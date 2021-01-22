package main

import (
	"crypto/sha1"
	"encoding/xml"
	"ffdaemon"
	"fflog"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"wxbizmsgcrypt"
)

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

func myGetSha1(strSrc string) string {
	t := sha1.New()
	io.WriteString(t, strSrc)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func modifyMachine(w http.ResponseWriter, r *http.Request )  {
	token := "RU0FtmLmkuZi834UGT76Dv"
	receiverId := "ww2a9e22ecf83c8f59"
	encodingAeskey := "ww6IXmYwcuGG6SbSFwk5iScpQ68NvVxagnRHiUzye71"
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)

	r.ParseForm()

	if r.Method == "POST" {
		reqMsgSign := r.FormValue("msg_signature")
		reqTimestamp := r.FormValue("timestamp")
		reqNonce := r.FormValue("nonce")
		reqData, _ := ioutil.ReadAll(r.Body)

		msg, cryptErr := wxcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, reqData)
		if nil != cryptErr {
			fflog.FFDebug("DecryptMsg fail" + cryptErr.ErrMsg)
			return
		}
		fflog.FFDebug("after decrypt msg: " + string(msg))

		var msgContent MsgContent
		err := xml.Unmarshal(msg, &msgContent)
		if nil != err {
			fflog.FFDebug("Unmarshal fail")
			return
		}

		var msgRsp MsgContent
		msgRsp.FromUsername = msgContent.ToUsername
		msgRsp.ToUsername = msgContent.FromUsername
		msgRsp.CreateTime = uint32(time.Now().Second())
		msgRsp.MsgType = msgContent.MsgType
		msgRsp.Agentid = msgContent.Agentid
		msgRsp.Msgid = msgContent.Msgid

		fflog.FFDebug("FromUsername:%s CreateTime:%d Content:%s",
			msgContent.FromUsername, msgContent.CreateTime, msgContent.Content)

		if msgContent.FromUsername != "stone" &&
			msgContent.FromUsername != "kevin" &&
			msgContent.FromUsername != "raymond" &&
			msgContent.FromUsername != "venson" &&
			msgContent.FromUsername != "rambo" &&
			msgContent.FromUsername != "william" &&
			msgContent.FromUsername != "Akai"{
			fflog.FFError(msgContent.FromUsername + " don't have permission")
			return
		}

		cmdStr := strings.Fields(msgContent.Content)
		if cmdStr[0] == "add" && len(cmdStr) == 3{
			cmd := exec.Command("/bin/bash", "-c",
				"/home/pf/code/run/authorizesvrd/cfg/allow_ip_add.sh " + cmdStr[1] + " " + cmdStr[2])
			strRet, err := cmd.Output()
			if err != nil{
				fflog.FFError("Do Shell Fail" + err.Error())
				return
			}
			msgRsp.Content = string(strRet)
		}else if cmdStr[0] == "del" && len(cmdStr) == 2{
			cmd := exec.Command("/bin/bash", "-c",
				"/home/pf/code/run/authorizesvrd/cfg/allow_ip_del.sh " + cmdStr[1])
			strRet, err := cmd.Output()
			if err != nil{
				fflog.FFError("Do Shell Fail" + err.Error())
				return
			}
			msgRsp.Content = string(strRet)
		}else if cmdStr[0] == "get" && len(cmdStr) == 2{
			cmd := exec.Command("/bin/bash", "-c",
				"/home/pf/code/run/authorizesvrd/cfg/allow_ip_get.sh " + cmdStr[1])
			strRet, err := cmd.Output()
			if err != nil{
				fflog.FFError("Do Shell Fail" + err.Error())
				return
			}
			msgRsp.Content = string(strRet)
		} else{
			msgRsp.Content = "CMD should like `add IP UUID` for add or like `del IP` for del"
		}

		retStr, err := xml.Marshal(msgRsp)
		encryptMsg, cryptErr := wxcpt.EncryptMsg(string(retStr), reqTimestamp, reqNonce)
		if cryptErr != nil{
			fflog.FFDebug("EncryptMsg fail" + cryptErr.ErrMsg)
			return
		}

		fflog.FFDebug("Send Rsp:%s", string(encryptMsg))
		w.Write(encryptMsg)
	}else{
		msgSignature := r.FormValue("msg_signature")
		timeStamp := r.FormValue("timestamp")
		nonce := r.FormValue("nonce")
		echoStr := r.FormValue("echostr")

		rspEchoStr, cryptErr := wxcpt.VerifyURL(msgSignature, timeStamp, nonce, echoStr)
		if nil != cryptErr {
			fmt.Println("verifyUrl fail", cryptErr)
			return
		}

		w.Write([]byte(rspEchoStr))
		fflog.FFDebug("Verify suc msgSignature=%s, timeStamp=%s, nonce=%s" +
			", echoStr=%s", msgSignature, timeStamp, nonce, echoStr)
	}
}

func main(){
	ffdaemon.Daemon()
	fflog.Open()
	defer fflog.Close()

	fflog.FFDebug("Listen on Port 16668")
	http.HandleFunc("/modifyMachine", modifyMachine)

	err := http.ListenAndServe(":16668", nil)
	if err != nil{
		fflog.FFError("Listen 16668 fail" + err.Error())
		return
	}
}
