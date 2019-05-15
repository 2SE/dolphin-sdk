package mock

import (
	"github.com/2se/dolphin-sdk/trace"
)

var tr = trace.GetTracer()

func SendRequest() {
	tr.GetTrace()
}
