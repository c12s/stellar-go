package main

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
)

func FromRequest(r *http.Request, name string) (*Span, error) {
	traceId := r.Context().Value(trace).(Values).Get(trace_id)[0]
	spanId := r.Context().Value(trace).(Values).Get(span_id)[0]
	tags := r.Context().Value(trace).(Values).Get(tags)[0] //k:v;kv;...;kv:kv
	if traceId != "" && spanId != "" {
		span := InitSpan(NewSpanContext(traceId, spanId), name)
		defer span.Start()

		if tags != "" {
			span.ingestTags(tags)
		}
		return span, nil
	}
	return nil, errors.New("No trace context in request")
}

func FromContext(ctx context.Context, name string) (*Span, error) {
	traceId := ctx.Value(trace).(*Values).Get(trace_id)[0]
	spanId := ctx.Value(trace).(*Values).Get(span_id)[0]
	tags := ctx.Value(trace).(*Values).Get(tags)[0] //k:v;kv;...;kv:kv
	if traceId != "" && spanId != "" {
		span := InitSpan(NewSpanContext(traceId, spanId), name)
		defer span.Start()

		if tags != "" {
			span.ingestTags(tags)
		}
		return span, nil
	}
	return nil, errors.New("No trace in context")
}

func FromGRPCContext(ctx context.Context, name string) (*Span, error) {
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		traceId := md[trace_id][0]
		spanId := md[span_id][0]
		tags := md[tags][0] //k:v;kv;...;kv:kv
		if traceId != "" && spanId != "" {
			span := InitSpan(NewSpanContext(traceId, spanId), name)
			defer span.Start()

			if tags != "" {
				span.ingestTags(tags)
			}
			return span, nil
		}
	}
	return nil, errors.New("No trace in context")
}

func NewTracedRequest(method, url string, body io.Reader, span Spanner) (*http.Request, error) {
	c := context.WithValue(context.Background(), trace, span.Serialize())
	return http.NewRequestWithContext(c, method, url, body)
}

func NewTracedContext(ctx context.Context, span Spanner) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, trace, span.Serialize())
	} else {
		return context.WithValue(context.Background(), trace, span.Serialize())
	}
}

func NewTracedGRPCContext(ctx context.Context, span Spanner) context.Context {
	if ctx != nil {
		return metadata.NewOutgoingContext(context.Background(), span.Serialize().md)
	} else {
		return metadata.NewOutgoingContext(context.Background(), span.Serialize().md)
	}
}
