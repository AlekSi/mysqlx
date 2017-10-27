package mysqlx

import (
	"log"
	"os"
	"sync"
)

// We keep trace functions in global map for two reasons:
// 1. We want to reuse logger instance for all connections with given trace prefix.
// 2. More importantly, we want to pass t.Logf as a trace function for much better logging during tests,
//    and there is no better way to do it.

// Trace function signature.
type TraceFunc func(format string, v ...interface{})

// noTrace is a trace functions which does nothing.
// TODO check it is inlined and eliminated by compiler.
func noTrace(string, ...interface{}) {}

var (
	traceFuncs   = make(map[string]TraceFunc)
	traceFuncsRW sync.RWMutex
)

// getTracef returns trace function of logger with given prefix, creating it if required.
func getTracef(prefix string) TraceFunc {
	traceFuncsRW.RLock()
	tracef := traceFuncs[prefix]
	traceFuncsRW.RUnlock()
	if tracef != nil {
		return tracef
	}

	traceFuncsRW.Lock()
	tracef = log.New(os.Stderr, prefix, log.Lshortfile).Printf
	traceFuncs[prefix] = tracef
	traceFuncsRW.Unlock()
	return tracef
}

// setTestTracef sets trace function. Used only in tests with t.Logf.
func setTestTracef(prefix string, tracef TraceFunc) {
	traceFuncsRW.Lock()
	traceFuncs[prefix] = tracef
	traceFuncsRW.Unlock()
}
