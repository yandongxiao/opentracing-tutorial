package main

import (
	"context"
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
	// The StartSpanFromContext(下面的两个函数都调用它来创建新的Span) function uses
	// opentracing.GlobalTracer() to start the new spans, so we need to initialize
	// that global variable to our instance of Jaeger tracer
	opentracing.SetGlobalTracer(tracer)

	helloTo := os.Args[1]
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)
	defer span.Finish()

	// 做一个转换
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
}

// The Context allows storing arbitrary key-value pairs, so we can use it to store the
// currently active span. The OpenTracing API integrates with context.Context and provides
// convenient helper functions.
func formatString(ctx context.Context, helloTo string) string {
	// Note that we ignore the second value returned by the function, which is another
	// instance of the Context with the new span stored in it.
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}
