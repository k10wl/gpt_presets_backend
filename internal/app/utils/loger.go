package utils

import "fmt"

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

func (l *Logger) Log(text string, args ...interface{}) {
	fmt.Printf("["+l.prefix+"]: "+text+"\n"+"%+v*\n", args...)
}
