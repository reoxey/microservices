package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		log.New(os.Stdout, "[account] ", log.Lshortfile| log.Lmsgprefix| log.LstdFlags),
	}
}

func (l *Logger) Metrics(next http.HandlerFunc, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer l.Println(kind, time.Since(t).String())
		next(w, r)
	}
}
