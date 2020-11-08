package  fflog

import (
    "bufio"
    "fmt"
    "os"
    "runtime"
    "sync"
    "time"
)

const(
	LOG_TRACE = 1
	LOG_DEBUG = 2
	LOG_INFO = 3
	LOG_WARNING = 4
	LOG_ERROR = 5
	LOG_CRITICAL = 6
)

type Log_Handle struct{
	log_level_  int8
	file_prefix_    string
	file_cur_   *os.File
    file_writer_    *bufio.Writer
    log_ticker_ *time.Ticker
    ch_         chan(bool)
    wg_          *sync.WaitGroup
}

var log_level_string = [...]string{"", "LM_TRACE", "LM_DEBUG", "LM_INFO","LM_WARN", "LM_ERROR", "LM_CRITCAL"}
var log *Log_Handle

func FFLog(log_level int8, v string, args ...interface{}){
    if log == nil {
        return
    }

    if log_level < log.log_level_ {
        return
    }

    str := getNowTimeMsec() + "@" + log_level_string[log_level] + "@" + fmt.Sprintf(v, args...)
    fmt.Println(str)
    log.file_writer_.WriteString(str + "\n")
}

func FFError(v string, args ...interface{}){
    FFLog(LOG_ERROR, v, args...)
}

func FFDebug(v string, args ...interface{}){
    FFLog(LOG_DEBUG, v, args...)
}

func FFInfo(v string, args ...interface{}){
    FFLog(LOG_INFO, v, args...)
}

func FFCrit(v string, args ...interface{}){
    FFLog(LOG_CRITICAL, v, args...)
}

func getNowTimeMsec() string{
        time_now := time.Now()
        time_fmt := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d %v", time_now.Year(),
            time_now.Month(), time_now.Day(), time_now.Hour(),
            time_now.Minute(), time_now.Second(), (time_now.UnixNano()/1e6 - time_now.Unix() * 1e3))
        return time_fmt
}

func Open() (*Log_Handle, error){
    _, file, _, _ := runtime.Caller(1)
    return OpenEx(file, LOG_DEBUG)
}

func OpenEx(prefix string, log_level int8) (*Log_Handle, error){
    if log != nil {
        return log, nil
    }

    wg := sync.WaitGroup{}
    return openWg(&wg, prefix, log_level)
}

func (log *Log_Handle) logTimer(){
    defer log.log_ticker_.Stop()
    for{
        select {
        case <- log.log_ticker_.C:
            log.file_writer_.WriteString(getNowTimeMsec() + "@" +
                log_level_string[log.log_level_] + "@" + "Log Flush!.\n")
            log.file_writer_.Flush()
        case stop := <-log.ch_:
            if stop{
                log.file_writer_.WriteString(getNowTimeMsec() + "@" +
                    log_level_string[log.log_level_] + "@" + "Log Stop!.\n")
                log.file_writer_.Flush()
                log.file_cur_.Close()
                log.wg_.Done()
                return
            }
        }
    }
}

func openWg(wg *sync.WaitGroup, prefix string, log_level int8) (*Log_Handle, error){
    log = new(Log_Handle)
    log.file_prefix_ = prefix
    log.log_level_ = log_level
    log.wg_ = wg

    time_now := time.Now()
    time_fmt := fmt.Sprintf("%d%02d%02d", time_now.Year(),
    time_now.Month(), time_now.Day())

    path := log.file_prefix_ + "_" + time_fmt
    var err error
    log.file_cur_, err = os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    log.log_ticker_ = time.NewTicker(time.Second)
    log.file_writer_ = bufio.NewWriter(log.file_cur_)
    log.ch_ = make(chan bool)

    wg.Add(1)
    go log.logTimer()

    return log, err
}

func Close(){
   log.ch_ <- true
   close(log.ch_)
   log.wg_.Wait()
}


