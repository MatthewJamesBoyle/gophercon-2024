package main

import (
	"context"
	"github.com/matthewjamesboyle/gophercon-2024/internal/flight"
	"github.com/matthewjamesboyle/gophercon-2024/internal/hotel"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/matthewjamesboyle/gophercon-2024/internal/transporthttp"
	"github.com/uhthomas/slogctx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var tracer trace2.Tracer

func main() {

	ctx := context.Background()

	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Create Trace Provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("app-one"),
		)),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("app-one")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx = slogctx.With(ctx, logger)

	c := &http.Client{
		Timeout: time.Second & 5,
	}

	ff := flight.NewService()
	hf := hotel.NewService(c, "http://localhost:8080")

	svc := recomendation.NewService(hf, ff)

	m := transporthttp.NewMux(ctx, svc)

	if err := http.ListenAndServe(":3000", m); err != nil {
		log.Fatal(err)
	}
}
