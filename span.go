package main

import (
	"fmt"
	sPb "github.com/c12s/scheme/stellar"
	"github.com/golang/protobuf/proto"
	"github.com/rs/xid"
	"strings"
	"time"
)

type Spanner interface {
	Child(name string) *Span
	AddLog(kv *KV)
	AddTag(kv *KV)
	AddBaggage(kv *KV)
	StartTime(startTime int64)
	EndTime(endtime int64)
	Finish() // send data to collecto and maybe serialize to ctx ot request
	Serialize() *Values
	Start()
	Marshall() ([]byte, error)
}

type SpanContext struct {
	s *sPb.SpanContext
}

func (s *SpanContext) String() string {
	return fmt.Sprintf("Ctx:(tid: %s sid: %s pid: %s) ", s.s.TraceId, s.s.SpanId, s.s.ParrentSpanId)
}

func NewSpanContext(tid, pid string) *SpanContext {
	return &SpanContext{
		&sPb.SpanContext{
			TraceId:       tid,
			SpanId:        xid.New().String(),
			ParrentSpanId: pid,
			Baggage:       map[string]string{},
		},
	}
}

type Span struct {
	s *sPb.Span
}

func InitSpan(c *SpanContext, name string) *Span {
	return &Span{
		&sPb.Span{
			SpanContext: c.s,
			Name:        name,
			Logs:        map[string]string{},
			Tags:        map[string]string{},
		},
	}
}

func (s *Span) Child(name string) *Span {
	context := NewSpanContext(s.s.SpanContext.TraceId, s.s.SpanContext.ParrentSpanId)
	span := InitSpan(context, name)
	defer span.Start()
	return span
}

func (s *Span) AddLog(kv *KV) {
	s.s.Logs[kv.key] = kv.value
}

func (s *Span) AddTag(kv *KV) {
	s.s.Tags[kv.key] = kv.value
}

func (s Span) AddBaggage(kv *KV) {
	s.s.SpanContext.Baggage[kv.key] = kv.value
}

func (s *Span) StartTime(t int64) {
	s.s.StartTime = t
}

func (s *Span) EndTime(t int64) {
	s.s.EndTime = t
}

func (s *Span) Finish() {
	s.s.EndTime = time.Now().Unix()
	//send to colledtor
	fmt.Println(fmt.Sprintf("Span.Finish() %d", (s.s.EndTime - s.s.StartTime)))
	data, err := s.Marshall()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)
}

func (span *Span) Serialize() *Values {
	s := map[string][]string{}
	s[trace_id] = []string{span.s.SpanContext.TraceId}
	s[span_id] = []string{span.s.SpanContext.SpanId}
	s[tags] = []string{span.digestTags()}
	return &Values{md: s}
}

func (s *Span) ingestTags(existing string) {
	for _, pair := range strings.Split(existing, tag_sep) {
		val := strings.Split(pair, pair_sep)
		s.AddTag(&KV{key: val[0], value: val[1]})
	}
}

func (s *Span) digestTags() string {
	t := []string{}
	for k, v := range s.s.Tags {
		t = append(t, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(t, tag_sep)
}

func (s *Span) String() string {
	return fmt.Sprintf("Span: %s %s", s.s.SpanContext, s.s.Name)
}

func (s *Span) Start() {
	s.s.StartTime = time.Now().Unix()
}

func (s *Span) Marshall() ([]byte, error) {
	return proto.Marshal(s.s)
}
