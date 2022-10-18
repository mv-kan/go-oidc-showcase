package main

import (
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/http/controller"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/http/router"
)

const (
	CheckTokenEndpoint = "CHECK_TOKEN_ENDPOINT"
	OPURL              = "OP_URL"

	LogFilePath = "LOG_FILEPATH"
)

func run() {
	// load env file
	err := godotenv.Load("op.env")
	if err != nil {
		log.Fatal(err.Error())
	}
	// configs
	opconf := config.OP{
		URL:                os.Getenv(OPURL),
		CheckTokenEndpoint: os.Getenv(CheckTokenEndpoint),
	}
	logconf := config.Logger{
		FilePath: os.Getenv(LogFilePath),
	}
	// create logger
	f, err := os.OpenFile(logconf.FilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	logger := log.New(map[log.LogOutput]io.Writer{
		log.FileOutput:    f,
		log.ConsoleOutput: os.Stdout,
	})
	// get controller
	messageCtrl := controller.NewProtectedMessage(logger, opconf)

	// get router
	r := router.New(messageCtrl)

	rsHost := os.Getenv("RS_HOST")
	logger.Info("run server on " + rsHost)
	http.ListenAndServe(rsHost, r)
}

func main() {
	run()
}
