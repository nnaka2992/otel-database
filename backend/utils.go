package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func httpError(w http.ResponseWriter, m string) {
	code, err := strconv.Atoi(m[:3])
	if err != nil {
		code = http.StatusInternalServerError
	}
	http.Error(w, m, code)
	log.Printf(m)
}

func readJson(r *http.Request) (map[string]interface{}, error) {
	l := r.ContentLength
	body := make([]byte, l)
	_, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
	}
	var params map[string]interface{}
	err = json.Unmarshal(body[:l], &params)
	if err != nil {
		return nil, fmt.Errorf("%d Bad Request: %s", http.StatusBadRequest, err)
	}
	return params, nil
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(
		stdout.WithPrettyPrint(),
		stdout.WithWriter(os.Stdout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stdout exporter: %w", err)
	}

	// Create a new tracer provider with the exporter and a sampler.
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register the trace provider with the global tracer provider.
	otel.SetTracerProvider(provider)

	// Set the global propagator to tracecontext so that the trace and span IDs from the incoming request
	// are extracted and propagated to the outgoing requests.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return provider, nil
}
