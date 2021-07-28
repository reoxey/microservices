// Author: reoxey
// Date: 22/07/2020 8:43 PM

package log

import (
	"fmt"
	l "log"
	"os"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var x *l.Logger

type GenericErr struct {
	Code int
	Err  error
}

func (c GenericErr) Error() string {
	return c.Error()
}

type ValidationErr struct {
	Err error
}

func (c ValidationErr) Error() string {
	return c.Error()
}

func init() {
	x = &l.Logger{}
	x.SetFlags(l.LstdFlags)
}

func Obj() *l.Logger {
	x.SetOutput(os.Stdout)
	return x
}

func Err(err error) error {
	prefix := trace(3)
	x.SetOutput(os.Stderr)
	x.Printf(prefix+" | %v", err)
	return err
}

func GrpcErr(err error) error {
	prefix := trace(3)
	x.SetOutput(os.Stderr)
	x.Printf(prefix+" | %v", err)
	return status.Error(codes.Internal, err.Error())
}

func GrpcCodeErr(err error, code codes.Code) error {
	prefix := trace(3)
	x.SetOutput(os.Stderr)
	x.Printf(prefix+" | %v", err)
	return status.Error(code, err.Error())
}

func Error(err error) {
	prefix := trace(3)
	x.SetOutput(os.Stderr)
	x.Printf(prefix+" | %v", err)
}

func Fatal(err error) {
	prefix := trace(3)
	x.SetOutput(os.Stderr)
	x.Fatalf(prefix+" | %v", err)
}

func Msg(s string) {
	prefix := trace(1)
	x.SetOutput(os.Stdout)
	x.Println(prefix, s)
}

func Debug(s string) {
	prefix := trace(3)
	x.SetOutput(os.Stdout)
	x.Println(prefix, s)
}

func trace(n int) string {
	var pre []string
	for i := 0; i < n; i++ {
		_, f, p, ok := runtime.Caller(i + 2)
		if ok {
			x := strings.Split(f, "/")
			pre = append(pre, fmt.Sprintf("%s:%d", x[len(x)-1], p))
		}
	}
	return strings.Join(pre, " ⮜ ")
}
