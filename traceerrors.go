package traceerrors

import (
	"fmt"
	"runtime"
	"strings"
)

const includeStackInError = true

// TraceError wraps an error with a message and a stack frame.
type TraceError struct {
	Msg   string
	Err   error
	Frame string
}

// Error implements the error interface.
func (e *TraceError) Error() string {
	var b strings.Builder
	if e.Msg != "" {
		b.WriteString(e.Msg)
		if e.Err != nil {
			b.WriteString(": ")
		}
	}
	if e.Err != nil {
		b.WriteString(e.Err.Error())
	}

	if includeStackInError && e.Frame != "" {
		b.WriteString("\n")
		b.WriteString(StackTrace(e))
	}

	return b.String()
}

// Unwrap returns the underlying error.
func (e *TraceError) Unwrap() error {
	return e.Err
}

// New creates a new TraceError with a message and a stack frame.
func New(msg string) error {
	return &TraceError{
		Msg:   msg,
		Frame: captureStackFrame(),
	}
}

// Newf creates a new TraceError with a formatted message and a stack frame.
func Newf(format string, args ...interface{}) error {
	return &TraceError{
		Msg:   fmt.Sprintf(format, args...),
		Frame: captureStackFrame(),
	}
}

// Wrap wraps an existing error with a message and a stack frame.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &TraceError{
		Msg:   msg,
		Err:   err,
		Frame: captureStackFrame(),
	}
}

// Wrapf wraps an existing error with a formatted message and a stack frame.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &TraceError{
		Msg:   fmt.Sprintf(format, args...),
		Err:   err,
		Frame: captureStackFrame(),
	}
}

// captureStackFrame captures the current stack frame.
func captureStackFrame() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	function := "unknown"
	if fn != nil {
		function = fn.Name()
	}
	return fmt.Sprintf("%s\n\t%s:%d", function, file, line)
}

// StackTrace returns the full stack trace by traversing the error chain.
func StackTrace(err error) string {
	var frames []string
	for err != nil {
		if te, ok := err.(*TraceError); ok {
			if te.Frame != "" {
				frames = append([]string{te.Frame}, frames...)
			}
			err = te.Err
		} else {
			break
		}
	}
	return strings.Join(frames, "\n")
}
