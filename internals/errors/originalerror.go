package errors

import (
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func getOriginalErrorStackTrace(err error) (stackTrace []string) {
	if err, ok := err.(stackTracer); ok {
		for _, frame := range err.StackTrace() {
			frameBs, err := frame.MarshalText()
			if err != nil {
				continue
			}
			stackTrace = append(stackTrace, string(frameBs))
		}
	}
	return stackTrace
}
