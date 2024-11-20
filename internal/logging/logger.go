package logging

import (
	"fmt"
	"os"
)

type Logger struct {
	out    *os.File
	errOut *os.File
}

func NewLogger() *Logger {
	return &Logger{
		out:    os.Stdout,
		errOut: os.Stderr,
	}
}

func (l *Logger) LogError(err error) {
	if _, outErr := l.errOut.WriteString(err.Error() + "\n"); outErr != nil {
		fmt.Fprintf(os.Stderr, "skipped logging entry due to error: %+v. Message: '%s'\n", outErr, err.Error())
	}
}

func (l *Logger) Log(msg string) {
	if _, err := l.out.WriteString(msg + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "skipped logging entry due to error: %+v. Message: '%s'\n", err, msg)
	}
}
