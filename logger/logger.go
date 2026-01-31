package logger

import (
	"CQS-KYC/config"
	"context"

	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger struct {
	*zap.Logger
}

func NewLogger(cfg *config.Config) *AppLogger {
	var coreArr []zapcore.Core
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.InfoLevel)
	coreArr = append(coreArr, consoleCore)
	log := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	return &AppLogger{
		Logger: log,
	}
}

func (l *AppLogger) InfoWithMask(ctx context.Context, message, secret string) {
	requestId := ctx.Value("requestId").(string)
	len := len(secret)
	if len > 0 && len <= 8 {
		secret = "*************"
	} else if len > 8 {
		secret = "*************" + secret[len-4:]
	}
	l.Logger.Info(message, zap.String("requestId", requestId), zap.String("data", secret))
}
