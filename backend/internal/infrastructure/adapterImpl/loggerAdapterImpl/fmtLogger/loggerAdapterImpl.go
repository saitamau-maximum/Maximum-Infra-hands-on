package fmtlogger

import (
	"fmt"

	"example.com/infrahandson/internal/interface/adapter"
)

// FmtLogger は、標準出力にログを出すLogger実装です。
type FmtLogger struct{}

// NewFmtLogger は、FmtLoggerのコンストラクタです。
func NewFmtLogger() adapter.LoggerAdapter {
	return &FmtLogger{}
}

func (l *FmtLogger) Info(msg string, args ...any) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
	fmt.Printf("\n")
}

func (l *FmtLogger) Warn(msg string, args ...any) {
	fmt.Printf("[WARN] "+msg+"\n", args...)
	fmt.Printf("\n")
}

func (l *FmtLogger) Error(msg string, args ...any) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
	fmt.Printf("\n")
}
