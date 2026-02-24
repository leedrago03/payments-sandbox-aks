package logging

import (
	"context"
	"net/http"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// InitTracing sets up the global OpenTelemetry propagator to use B3 (Zipkin/Istio style).
// This function should be called once at service startup.
func InitTracing() {
	// Configure the global propagator to use B3 Multi Header format.
	// This ensures that when we inject/extract context, we look for X-B3-TraceId, etc.
	p := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	otel.SetTextMapPropagator(p)
}

// ExtractTraceContext reads tracing headers from an incoming HTTP request
// and returns a new context containing the trace span.
func ExtractTraceContext(req *http.Request) context.Context {
	// 1. Get the global propagator (B3)
	propagator := otel.GetTextMapPropagator()
	
	// 2. Extract the context from the request headers
	ctx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))
	
	return ctx
}

// FastHTTPHeaderCarrier is a carrier for fasthttp.RequestHeader.
type FastHTTPHeaderCarrier struct {
	Header *fasthttp.RequestHeader
}

func (c FastHTTPHeaderCarrier) Get(key string) string {
	return string(c.Header.Peek(key))
}

func (c FastHTTPHeaderCarrier) Set(key string, value string) {
	c.Header.Set(key, value)
}

func (c FastHTTPHeaderCarrier) Keys() []string {
	keys := make([]string, 0, c.Header.Len())
	c.Header.VisitAll(func(key, value []byte) {
		keys = append(keys, string(key))
	})
	return keys
}

// ExtractTraceFromFastHTTP reads tracing headers from an incoming fasthttp request.
func ExtractTraceFromFastHTTP(ctx context.Context, header *fasthttp.RequestHeader) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, FastHTTPHeaderCarrier{Header: header})
}

// InjectTraceToFastHTTP injects the current trace span into fasthttp headers.
func InjectTraceToFastHTTP(ctx context.Context, header *fasthttp.RequestHeader) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, FastHTTPHeaderCarrier{Header: header})
}

// InjectTraceContext injects the current trace span from the context
// into the headers of an outgoing HTTP request.
func InjectTraceContext(ctx context.Context, req *http.Request) {
	// 1. Get the global propagator (B3)
	propagator := otel.GetTextMapPropagator()
	
	// 2. Inject the context into the request headers
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
}

// TraceAwareClient is an HTTP client wrapper that automatically injects trace context.
type TraceAwareClient struct {
	Client *http.Client
}

// NewTraceAwareClient returns a new TraceAwareClient wrapping a standard http.Client.
func NewTraceAwareClient(client *http.Client) *TraceAwareClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &TraceAwareClient{Client: client}
}

// Do performs an HTTP request, injecting the trace context from req.Context() into the headers.
func (c *TraceAwareClient) Do(req *http.Request) (*http.Response, error) {
	// Inject trace context before sending
	InjectTraceContext(req.Context(), req)
	return c.Client.Do(req)
}
