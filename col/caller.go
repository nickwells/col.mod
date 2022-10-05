package col

import (
	"fmt"
	"runtime"
)

// caller returns a string giving the filename and line number of the caller
// of the calling function. This is intended for providing useful debugging
// messages. Note that we ask for the second stack entry above this: 0 would
// give the location of the call to runtime.Caller, 1 would give the location
// of the call to caller() but we want to see where the parent function was
// called so we pass 2
func caller() string {
	if pc, file, line, ok := runtime.Caller(2); ok {
		f := runtime.FuncForPC(pc)
		funcName := "unknown"
		if f != nil {
			funcName = f.Name()
		}
		return fmt.Sprintf("%s:%d [%s]", file, line, funcName)
	}
	return "unknown-file:0 [unknown]"
}
