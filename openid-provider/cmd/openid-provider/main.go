package main

import (
	"io"
	"os"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
)

func main() {
	// test logging
	log.Info("Testing message")
	f, err := os.OpenFile("./test.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	logger := log.New(map[log.LogOutput]io.Writer{
		log.FileOutput: f,
	})
	logger.Info("file log testing message")
}
