package main
import (
	"fflog"
)

func main() {
	fflog.Open()
	fflog.FFLog(fflog.LOG_DEBUG, "Will Give Right")
	fflog.Close()
}
