package transporthttp

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VictoriaMetrics/metrics"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/uhthomas/slogctx"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"time"
)

var (
	requestDuration = metrics.NewSummary("request_duration_seconds")
)

func NewMux(ctx context.Context, svc *recomendation.Service) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/trip/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			recordDuration(startTime)
		}()

		// Extract budget from query parameters and parse it
		budget := r.URL.Query().Get("budget")
		b, err := strconv.ParseInt(budget, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tr := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("gophercon-2024")
		ctx, span := tr.Start(r.Context(), "get_recommendation")
		defer span.End()

		rec, err := svc.Get(ctx, int(b))
		if err != nil {
			switch {
			case errors.Is(err, recomendation.ErrBudgetOutOfBounds):
				w.WriteHeader(http.StatusBadRequest)
			default:
				slogctx.From(ctx).Error("unhandled error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		span.End()

		// Marshal response
		res, err := json.Marshal(rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	})

	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func recordDuration(startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	requestDuration.Update(duration)
}
