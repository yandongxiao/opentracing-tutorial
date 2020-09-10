package main

import (
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go/log"

	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

// 1. creates a trace that consists of a single span.
// 2. the single span combined two operations performed by the program, formatting the output string and printing it.
func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	// new tracer
	// a tracer instance is used to start new spans via StartSpan function
	// It(service) is used to mark all spans emitted by the tracer as originating from a hello-world service.
	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()

	// each span is given an operation name, "say-hello" in this case
	// 关于operationName
	//  1. the operation name is meant to represent a class of spans, rather than a unique instance.
	//  2. 你不应该将参数名称（一个不断变化的值）传递给StartSpan. 解决办法如下：
	//  3. The recommended solution is to annotate spans with tags or logs.
	helloTo := os.Args[1]
	span := tracer.StartSpan("say-hello")

	// A tag is a key-value pair that provides certain metadata about the span.
	// The tags are meant to describe attributes of the span that apply to the whole duration of the span.
	span.SetTag("hello-to", helloTo)

	// A log is similar to a regular log statement, it contains a timestamp and some data,
	// but it is associated with span from which it was logged.
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		// Just google "structured-logging" for many articles on this topic.
		// structured logging APIs encourage you to separate bits and pieces of that message into key-value
		// pairs that can be automatically processed by log aggregation systems.
		// The OpenTracing Specification also recommends all log statements to contain an event field that
		// describes the overall event being logged, with other attributes of the event provided as additional fields.
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	println(helloStr)
	span.LogKV("event", "println")

	// each span must be finished by calling its Finish() function
	// the start and end timestamps of the span will be captured automatically by the tracer implementation
	span.Finish()
}