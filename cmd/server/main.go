package main

import (
	"context"
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

	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close() // Ensure the file is closed when the program exits

	ctx := context.Background()
	// Create a logger that writes to the log file
	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	ctx = slogctx.With(ctx, logger)

	c := &http.Client{
		Timeout: time.Second & 5,
	}

	ff := flight.NewService()
	hf := hotel.NewService(c, "http://localhost:8080")

	svc := recomendation.NewService(hf, ff)

	m := transporthttp.NewMux(ctx, svc)

	logger.InfoContext(ctx, "app_starting")
	if err := http.ListenAndServe(":3000", m); err != nil {
		log.Fatal(err)
	}
}
