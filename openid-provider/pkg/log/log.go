package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logging by instance
type LogOutput int

const (
	FileOutput LogOutput = iota
	ConsoleOutput
)

// to add console writer you need to add this value in map
// ConsoleOutput : os.Stdout
// to add file log writer add this in map
// FileOutput : io.Writer (this is open io writer to file)
func New(outputmap map[LogOutput]io.Writer) Logger {
	zapLogger := newZapLogger(outputmap)
	return logger{zapLogger}
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

type logger struct {
	zapLogger *zap.Logger
}

func (l logger) Debug(msg string) {
	l.zapLogger.Debug(msg)
}
func (l logger) Info(msg string) {
	l.zapLogger.Info(msg)
}
func (l logger) Error(msg string) {
	l.zapLogger.Error(msg)
}
func (l logger) Fatal(msg string) {
	l.zapLogger.Fatal(msg)
}

// logging on global level
var (
	localLogger = newZapLogger(map[LogOutput]io.Writer{
		ConsoleOutput: os.Stdout,
	})
)

func Debug(msg string) {
	localLogger.Debug(msg)
}

func Info(msg string) {
	localLogger.Info(msg)
}

func Error(msg string) {
	localLogger.Error(msg)
}

func Fatal(msg string) {
	localLogger.Fatal(msg)
}

// to add console writer you need to add this value in map
// ConsoleOutput : os.Stdout
// to add file log writer add this in map
// FileOutput : io.Writer (this is open io writer to file)
func newZapLogger(outputmap map[LogOutput]io.Writer) *zap.Logger {
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultLogLevel := zapcore.DebugLevel

	cores := make([]zapcore.Core, 0)
	for output, stream := range outputmap {
		switch output {
		case ConsoleOutput:
			consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)
			core := zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(stream), defaultLogLevel),
			)
			cores = append(cores, core)
		case FileOutput:
			fileEncoder := zapcore.NewJSONEncoder(zapConfig)
			core := zapcore.NewTee(
				zapcore.NewCore(fileEncoder, zapcore.AddSync(stream), defaultLogLevel),
			)
			cores = append(cores, core)
		}
	}
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.InfoLevel))
	return zapLogger
}

// func newZapLogger() *zap.Logger {
// 	zapConfig := zap.NewProductionEncoderConfig()
// 	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

// 	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)

// 	defaultLogLevel := zapcore.DebugLevel

// 	core := zapcore.NewTee(
// 		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
// 	)

// 	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.InfoLevel))

// 	return zapLogger
// }
