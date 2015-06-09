package main

import (
	"flag"
	"log"
	"os"
)

var skinDir = ""

func main() {
	listen := flag.String("l", "127.0.0.1:8080", "address to listen on")
	skinDirFlag := flag.String("s", "/tmp/jwpack", "directory for storing skin files")

	flag.Parse()

	skinDir = *skinDirFlag

	err := os.MkdirAll(skinDir, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	initTemplates()
	initHandlers()

	listenAndServe(*listen)
}
