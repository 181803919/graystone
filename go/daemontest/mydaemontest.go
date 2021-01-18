package main

import (
	"ffdaemon"
	"fmt"
)

func main() {
	ffdaemon.Daemon()
	fmt.Println("I can See You")
}
