package biggo_dbg

import (
	"bytes"
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

	pParams := bytes.NewBuffer([]byte{})
	if len(params) > 0 {
		json.NewEncoder(pParams).Encode(params)
	}

	fmt.Fprintf(os.Stdout, "Call: %s%s\r\n", filepath.Base(f.Name()), pParams)
}

func SlowTrace(timeout uint8, params ...interface{}) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	pParams := bytes.NewBuffer([]byte{})
	if len(params) > 0 {
		json.NewEncoder(pParams).Encode(params)
	}

	fmt.Fprintf(os.Stdout, "Call: %s%s\r\n", filepath.Base(f.Name()), pParams)
	time.Sleep(time.Second * time.Duration(timeout))
}
