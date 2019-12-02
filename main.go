package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello world")
	trace := Init("test-trace")
	defer trace.Finish()

	span := trace.Span("first-span")
	defer span.Finish()

	fmt.Println(span)

	time.Sleep(2 * time.Second)

	child(NewTracedContext(context.TODO(), span))
	child2(NewTracedContext(context.TODO(), span))
}

func child(ctx context.Context) {
	child, _ := FromContext(ctx, "child-span")
	defer child.Finish()
	fmt.Println(child)
	cchild(NewTracedContext(ctx, child))
	time.Sleep(1 * time.Second)
}

func child2(ctx context.Context) {
	child, _ := FromContext(ctx, "child-span")
	defer child.Finish()
	fmt.Println(child)
	time.Sleep(1 * time.Second)
}

func cchild(ctx context.Context) {
	child, _ := FromContext(ctx, "child-child-span")
	defer child.Finish()
	fmt.Println(child)
	time.Sleep(2 * time.Second)
}
