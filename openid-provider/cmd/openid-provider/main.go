package main

import (
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
)

const (
	AuthEndpoint       = "AUTH_ENDPOINT"
	LoginEndpoint      = "LOGIN_ENDPOINT"
	TokenEndpoint      = "TOKEN_ENDPOINT"
	CheckTokenEndpoint = "CHECK_TOKEN_ENDPOINT"
	OPHost             = "OP_HOST"
	OPURL              = "OP_URL"

	LogFilePath = "LOG_FILEPATH"
)

func run() {
	// load env file
	err := godotenv.Load("op.env")
	if err != nil {
		log.Fatal(err.Error())
	}
	// set variables config
	httpconf := config.HTTP{
		AuthEndpoint:       os.Getenv(AuthEndpoint),
		LoginEndpoint:      os.Getenv(LoginEndpoint),
		TokenEndpoint:      os.Getenv(TokenEndpoint),
		CheckTokenEndpoint: os.Getenv(CheckTokenEndpoint),
		OPHost:             os.Getenv(OPHost),
		OPURL:              os.Getenv(OPURL),
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
	// get storage
	// get controllers

	logger.Info("file and console log testing message")
}

func main() {
	// test logging
	log.Info("Testing message")
	f, err := os.OpenFile("./test.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	logger := log.New(map[log.LogOutput]io.Writer{
		log.FileOutput:    f,
		log.ConsoleOutput: os.Stdout,
	})
	logger.Info("file and console log testing message")
}
