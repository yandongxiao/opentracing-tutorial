package main

import (
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()

	helloTo := os.Args[1]
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)
	defer span.Finish()

	helloStr := formatString(span, helloTo)
	printHello(span, helloStr)
}

// wrap each function into its own span.
// 这种封装的问题在于，你需要在函数之间，像传递Context一样传递Span。这对于开发来讲是很大的负担！
func formatString(rootSpan opentracing.Span, helloTo string) string {
	// If we think of the trace as a directed acyclic graph where nodes are the
	// spans and edges are the causal relationships between them, then the ChildOf
	// option is used to create one such edge between span and rootSpan
	//
	// ChildOf的返回值：In the API the edges are represented by SpanReference type
	// that consists of a SpanContext and a label.
	// ChildOf relationship means that the rootSpan has a logical dependency on the
	// child span before rootSpan can complete its operation
	// FollowsFrom, which means the rootSpan is the ancestor in the DAG, but it does
	// not depend on the completion of the child span
	span := rootSpan.Tracer().StartSpan(
		"formatString", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	return helloStr
}

// wrap each function into its own span.
func printHello(rootSpan opentracing.Span, helloStr string) {
	span := rootSpan.Tracer().StartSpan(
		"printHello", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()
	println(helloStr)
	span.LogKV("event", "println")
}
