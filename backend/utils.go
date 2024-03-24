package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
    "go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
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
	// OTLP exporter config for Collector (using default config)
	exporter, err := texporter.New(texporter.WithProjectID("otel-database"))
	if err != nil {
		return nil, err
	}
	res, err := resource.New(
		context.Background(),
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("otel-database"),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("production"),
			semconv.TelemetrySDKNameKey.String("opentelemetry"),
			semconv.TelemetrySDKLanguageKey.String("go"),
			semconv.TelemetrySDKVersionKey.String("1.24.0"),
		),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}
