package log

import (
	"fmt"
	"os"
	"time"
)

type Tracer struct {
	TraceId string
	Logger
}


func (t *Tracer) Debugf(format string, args ...interface{}) {
	t.Logger.Debugf(t.wrapTraceId()+format, args...)
}


func (t *Tracer) Infof(format string, args ...interface{}) {
	t.Logger.Infof(t.wrapTraceId()+format, args...)
}

func (t *Tracer) Warnf(format string, args ...interface{}) {
	t.Logger.Warnf(t.wrapTraceId()+format, args...)
}

func (t *Tracer) Errorf(format string, args ...interface{}) {
	t.Logger.Errorf(t.wrapTraceId()+format, args...)
}

func (t *Tracer) Fatalf(format string, args ...interface{}) {
	t.Logger.Fatalf(t.wrapTraceId()+format, args...)
	os.Exit(1)
}

func (t *Tracer) wrapTraceId() string {
	return "traceId:" + t.TraceId + " "
}

func Trace(traceId string) Logger {
	if traceId == "" {
		traceId = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return &Tracer{
		TraceId: traceId,
		Logger: DefaultLogger,
	}
}
