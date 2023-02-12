package util

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func Fatal(v ...any) {
	fmt.Fprintln(os.Stderr, time.Now().Local())
	pc, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s:%d %s()\n", file, line, runtime.FuncForPC(pc).Name())
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
