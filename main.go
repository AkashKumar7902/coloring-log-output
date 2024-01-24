package main

import (
	"bytes"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type colorConsoleEncoder struct {
	*zapcore.EncoderConfig
	zapcore.Encoder
}

func NewColorConsole(cfg zapcore.EncoderConfig) (enc zapcore.Encoder) {
	return colorConsoleEncoder{
		EncoderConfig: &cfg,
		// Using the default ConsoleEncoder can avoid rewriting interfaces such as ObjectEncoder
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

// EncodeEntry overrides ConsoleEncoder's EncodeEntry
func (c colorConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (buf *buffer.Buffer, err error) {
	buff, err := c.Encoder.EncodeEntry(ent, fields) // Utilize the existing implementation of zap
	if err != nil {
		return nil, err
	}

	bytesArr := bytes.Replace(buff.Bytes(), []byte("\\u001b"), []byte("\u001b"), -1)
	buff.Reset()
	buff.AppendString(string(bytesArr))
	return buff, err
}

func init() {
	_ = zap.RegisterEncoder("colorConsole", func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewColorConsole(config), nil
	})
}

func setupLogger() *zap.Logger {
	logCfg := zap.NewDevelopmentConfig()

	logCfg.Encoding = "colorConsole"

	logger, _ := logCfg.Build()
	return logger
}

func main() {
	ColorLogger := setupLogger()
	var HighlightGreen = color.New(color.FgGreen).SprintFunc()
	var HighlightYellow = color.New(color.FgYellow).SprintFunc()
	ColorLogger.Info("test log", zap.String(HighlightGreen("key"), HighlightYellow("value")))
}
