package main

import (
	"fmt"
	"github.com/rs/xid"
)

type Tracer interface {
	Span(name string) Spanner
	Finish()
}

type Trace struct {
	id   string
	name string
}

func (t *Trace) Span(name string) *Span {
	ctx := &SpanContext{
		traceid:   t.id,
		spanid:    xid.New().String(),
		parrentId: "-",
		baggage:   map[string]string{},
	}

	s := &Span{
		context: ctx,
		name:    name,
	}
	defer s.Start()
	return s
}

func (t *Trace) Finish() {
	//send to colledtor
	fmt.Println("Trace.Finish()")
}

func (t *Trace) String() string {
	return fmt.Sprintf("Trace: %s %s", t.id, t.name)
}

func Init(name string) *Trace {
	return &Trace{
		id:   xid.New().String(),
		name: name,
	}
}
