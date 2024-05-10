package biggo_dbg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func Trace(params ...interface{}) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	var pParams string
	if len(params) > 0 {
		if buffer, err := json.Marshal(struct{ Values []interface{} }{Values: params}); err == nil {
			pParams = fmt.Sprintf(" Params: %s", string(buffer))
		}
	}

	fmt.Fprintf(os.Stdout, "Call: %s%s\r\n", filepath.Base(f.Name()), pParams)
}

func SlowTrace(timeout uint8, params ...interface{}) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	var pParams string
	if len(params) > 0 {
		if buffer, err := json.Marshal(struct{ Values []interface{} }{Values: params}); err == nil {
			pParams = fmt.Sprintf(" Params: %s", string(buffer))
		}
	}

	fmt.Fprintf(os.Stdout, "Call: %s%s", filepath.Base(f.Name()), pParams)
	time.Sleep(time.Second * time.Duration(timeout))
}
