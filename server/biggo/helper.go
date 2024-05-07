package biggo

import (
	"encoding/json"
	"fmt"
	"os"
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

	fmt.Fprintf(os.Stdout, "=====TRACE===== [%s] Call: %s%s\r\n", time.Now(), f.Name(), pParams)
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

	fmt.Fprintf(os.Stdout, "===============\r\n=====SLOW1===== [%s] Call: %s%s\r\n===============\r\n", time.Now(), f.Name(), pParams)
	time.Sleep(time.Second * time.Duration(timeout))
}
