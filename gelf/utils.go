package gelf

import (
	"runtime"
	"strings"
)

const (
	pathVendor = `/vendor/`
	pathSrc    = `/src/`
	pathMod    = `/mod/`
)

// remove info of source location from path
func srcFile(path string) string {
	i := strings.LastIndex(path, pathVendor)
	if i >= 0 {
		return path[i+len(pathVendor):]
	}
	i = strings.LastIndex(path, pathSrc)
	if i >= 0 {
		return path[i+len(pathSrc):]
	}
	i = strings.LastIndex(path, pathMod)
	if i >= 0 {
		return path[i+len(pathMod):]
	}
	return path
}

// getCaller returns the filename and the line info of a function
// further down in the call stack.  Passing 0 in as callDepth would
// return info on the function calling getCallerIgnoringLog, 1 the
// parent function, and so on.  Any suffixes passed to getCaller are
// path fragments like "/pkg/log/log.go", and functions in the call
// stack from that file are ignored.
func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}
		file = srcFile(file)

		for _, s := range suffixesToIgnore {
			if strings.HasSuffix(file, s) {
				callDepth++
				continue outer
			}
		}
		break
	}
	return
}
