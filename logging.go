package logging

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	Debug   *Logger
	Info    *Logger
	Warning *Logger
	Error   *Logger

	setup   sync.Once
	base    *log.Logger
	discard *log.Logger
)

func doSetup() {
	base = log.New(os.Stderr, "", log.LstdFlags)
	discard = log.New(ioutil.Discard, "", log.LstdFlags)

	Debug = NewLogger("DEBUG", discard)
	Info = NewLogger("INFO", base)
	Warning = NewLogger("WARNING", base)
	Error = NewLogger("ERROR", base)
}

func Setup() {
	setup.Do(doSetup)
}

func Debugging() {
	Debug = NewLogger("DEBUG", base)
}
