package main

import "log"

func main() {
}

func init() {
	// LogFlag 日志显示格式
	lf := log.Ldate | log.Ltime | log.Lshortfile
	log.SetFlags(lf)
}
