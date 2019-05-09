package server

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
)

const (
	//invoke function trace
	baseCalen = 7
	//skip step to invoke func,it is very important
	invokeSkipForInit = 2
	invokeSkipForPush = 1
)

var tc *traceCache

func init() {
	initTraceCache()
	tc.rloop()
}

func GetTrace() string {
	return tc.getTrace()
}

func initTraceCache() {
	tc = &traceCache{
		tmap:      make(map[uintptr]string),
		receiveCh: make(chan *trace, 50),
		removeCh:  make(chan uintptr, 50),
		callers:   make(map[string]int),
	}
}

type trace struct {
	ptr     uintptr
	traceId string
}

type traceCache struct {
	//invoke func name
	curName string
	//the func before invoke func
	callers map[string]int
	ready   bool

	tmap      map[uintptr]string
	receiveCh chan *trace
	removeCh  chan uintptr
	m         sync.Mutex
}

//init invoke function trace
func (t *traceCache) initCurIndex() {
	t.m.Lock()
	defer t.m.Unlock()
	if t.ready {
		return
	}

	i := 1
	ptr, _, _, _ := runtime.Caller(invokeSkipForInit)
	t.curName = runtime.FuncForPC(ptr).Name()
	logrus.Info("runtime invoke ", t.curName)
Loop:
	ps := make([]uintptr, baseCalen*i)
	idx := runtime.Callers(invokeSkipForInit, ps)
	if idx == baseCalen*i {
		i++
		goto Loop
	}
	for k, v := range ps {
		f := runtime.FuncForPC(v)
		fnam := f.Name()
		if fnam != "" {
			t.callers[f.Name()] = k
		} else {
			break
		}
	}
	t.ready = true
}

//GetTrace/sendGrpc/.../invoke min:4
func (t *traceCache) getTrace() string {
	skip := 6
Loop:
	//fmt.Println(i)
	ptr, _, _, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(ptr)
	fnam := fn.Name()
	if fnam != t.curName {
		idx, ok := t.callers[fnam]
		if !ok {
			//not arrive
			skip = skip + baseCalen

		} else {
			skip = skip - idx - 1
		}
		goto Loop
	}
	return t.tmap[fn.Entry()]
}

func (t *traceCache) get(ptr uintptr) string {
	return t.tmap[ptr]
}

func (t *traceCache) push(traceId string) {
	if !t.ready {
		t.initCurIndex()
	}
	pc, _, _, _ := runtime.Caller(invokeSkipForPush)
	ptr := runtime.FuncForPC(pc).Entry()
	t.receiveCh <- &trace{ptr, traceId}
}
func (t *traceCache) remove() {
	pc, _, _, _ := runtime.Caller(2)
	ptr := runtime.FuncForPC(pc).Entry()
	t.removeCh <- ptr
}

func (t *traceCache) rloop() {
	go func() {
		for ch := range t.receiveCh {
			t.tmap[ch.ptr] = ch.traceId
		}
	}()
	go func() {
		for ch := range t.removeCh {
			delete(t.tmap, ch)
		}
	}()
}
