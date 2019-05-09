package mock

import (
	"fmt"
	"github.com/2se/dolphin-sdk/trace"
)

var tr = trace.GetTracer()

func SendRequest() {
	fmt.Println("mock traceId=>", tr.GetTrace())
}
