package logging

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

type Logger struct {
	prefix string
	Logger *log.Logger
}

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(3)

	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", path.Base(file), line)
}

func addObject(obj interface{}, arr []interface{}) []interface{} {
	return append([]interface{}{getFileAndLine(), obj}, arr...)
}

func prefix(prefix, format string) string {
	return getFileAndLine() + " " + prefix + " " + format
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Logger.Fatal(addObject(l.prefix, v)...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(prefix(l.prefix, format), v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Logger.Fatalln(addObject(l.prefix, v)...)
}

func (l *Logger) Flags() int {
	return l.Logger.Flags()
}

func (l *Logger) Output(calldepth int, s string) error {
	return l.Logger.Output(calldepth-1, s)
}

func (l *Logger) Panic(v ...interface{}) {
	l.Logger.Panic(addObject(l.prefix, v)...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(prefix(l.prefix, format), v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.Logger.Panicln(addObject(l.prefix, v)...)
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) Print(v ...interface{}) {
	l.Logger.Print(addObject(l.prefix, v)...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(prefix(l.prefix, format), v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.Logger.Println(addObject(l.prefix, v)...)
}

func (l *Logger) SetFlags(flag int) {
	l.Logger.SetFlags(flag)
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func NewLogger(prefix string, logger *log.Logger) *Logger {
	lo := new(Logger)
	lo.Logger = logger
	lo.prefix = prefix
	return lo
}
