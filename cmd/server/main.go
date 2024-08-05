package main

import (
	"context"
	"errors"
	"github.com/matthewjamesboyle/gophercon-2024/internal/flight"
	"github.com/matthewjamesboyle/gophercon-2024/internal/hotel"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/matthewjamesboyle/gophercon-2024/internal/transporthttp"
	"github.com/uhthomas/slogctx"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx := context.Background()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

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
