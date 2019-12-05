package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// trace := Init("test-trace")
	// defer trace.Finish()

	// span := trace.Span("first-span")
	// defer span.Finish()

	// fmt.Println(span)

	// child(NewTracedContext(context.Background(), span))
	// child2(NewTracedContext(context.Background(), span))

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(62*time.Second))
	defer cancel()

	n, err := NewCollector("0.0.0.0:4222", "collector")
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := InitCollector("logs/", n)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Start(ctx, 15*time.Second)
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
