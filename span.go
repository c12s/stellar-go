package main

import (
	"fmt"
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
}

type SpanContext struct {
	traceid   string
	spanid    string
	parrentId string
	baggage   map[string]string // propagated to other children and on other spans
}

func (s *SpanContext) String() string {
	return fmt.Sprintf("Ctx:(tid: %s sid: %s pid: %s) ", s.traceid, s.spanid, s.parrentId)
}

type Span struct {
	context   *SpanContext
	name      string
	logs      map[string]string
	tags      map[string]string // propagated to other children and on other spans
	startTime int64
	endTime   int64
}

func (s *Span) Child(name string) *Span {
	ctx := &SpanContext{
		traceid:   s.context.traceid,
		spanid:    xid.New().String(),
		parrentId: s.context.spanid,
		baggage:   map[string]string{},
	}

	span := &Span{
		context: ctx,
		name:    name,
	}
	defer span.Start()
	return span
}

func (s *Span) AddLog(kv *KV) {
	s.logs[kv.key] = kv.value
}

func (s *Span) AddTag(kv *KV) {
	s.tags[kv.key] = kv.value
}

func (s Span) AddBaggage(kv *KV) {
	s.context.baggage[kv.key] = kv.value
}

func (s *Span) StartTime(t int64) {
	s.startTime = t
}

func (s *Span) EndTime(t int64) {
	s.endTime = t
}

func (s *Span) Finish() {
	s.endTime = time.Now().Unix()
	//send to colledtor
	fmt.Println(fmt.Sprintf("Span.Finish() %d", (s.endTime - s.startTime)))
}

func (span *Span) Serialize() *Values {
	s := map[string][]string{}
	s[trace_id] = []string{span.context.traceid}
	s[span_id] = []string{span.context.spanid}
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
	for k, v := range s.tags {
		t = append(t, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(t, tag_sep)
}

func (s *Span) String() string {
	return fmt.Sprintf("Span: %s %s", s.context, s.name)
}

func (s *Span) Start() {
	s.startTime = time.Now().Unix()
}
