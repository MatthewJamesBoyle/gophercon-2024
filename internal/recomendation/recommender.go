package recomendation

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type finder interface {
	Find(ctx context.Context, budget int) (string, error)
}

func NewService(hotelFinder finder, flightFinder finder) *Service {
	return &Service{hotelFinder: hotelFinder, flightFinder: flightFinder}
}

type Service struct {
	hotelFinder  finder
	flightFinder finder
}

type Recommendation struct {
	Hotel  string
	Flight string
}

var ErrBudgetOutOfBounds = errors.New("budget out of bounds")

func (svc *Service) Get(ctx context.Context, budgetInPounds int) (*Recommendation, error) {
	if budgetInPounds == 0 {
		return nil, fmt.Errorf("must be greater than 0: %w", ErrBudgetOutOfBounds)
	}
	if budgetInPounds > 100000 {
		return nil, fmt.Errorf("wow bid spender, use a different API: %w", ErrBudgetOutOfBounds)
	}

	// use errGroup to make this faster
	tr := trace.SpanFromContext(ctx).TracerProvider().Tracer("gophercon-2024")
	ctx, hotelSpan := tr.Start(ctx, "get_hotel")
	hotel, err := svc.hotelFinder.Find(ctx, budgetInPounds)
	if err != nil {
		return nil, fmt.Errorf("failed to find hotel in budget: %w", err)
	}
	hotelSpan.AddEvent("hotel_selected", trace.WithAttributes(attribute.String("hotel", hotel)))
	hotelSpan.End()

	ctx, flightSpan := tr.Start(ctx, "get_flight")
	flight, err := svc.flightFinder.Find(ctx, budgetInPounds)
	if err != nil {
		return nil, fmt.Errorf("failed to find flight in budget: %w", err)
	}
	flightSpan.End()

	return &Recommendation{
		Hotel:  hotel,
		Flight: flight,
	}, nil
}
