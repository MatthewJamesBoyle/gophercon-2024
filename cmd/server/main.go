package main

import (
	"context"
	"github.com/matthewjamesboyle/gophercon-2024/internal/flight"
	"github.com/matthewjamesboyle/gophercon-2024/internal/hotel"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/matthewjamesboyle/gophercon-2024/internal/transporthttp"
	"log"
	"net/http"
	"time"
)

func main() {

	ctx := context.Background()

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
