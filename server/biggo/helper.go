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
			pParams = string(buffer)
		}
	}

	fmt.Fprintf(os.Stdout, "[%s] Call: %s Params: %s\r\n", time.Now(), f.Name(), pParams)
}
