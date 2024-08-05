package transporthttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/uhthomas/slogctx"
	"net/http"
	"strconv"
	"time"
)

var (
	requestsTotal   = "requests_total{status_code=%d}"
	requestDuration = metrics.NewSummary("request_duration_seconds")
)

func NewMux(ctx context.Context, svc *recomendation.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/trip/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer recordDuration(startTime)

		// Extract budget from query parameters and parse it
		budget := r.URL.Query().Get("budget")
		b, err := strconv.ParseInt(budget, 10, 64)
		if err != nil {
			httpError(ctx, w, http.StatusBadRequest, err)
			return
		}

		// Get recommendation
		rec, err := svc.Get(ctx, int(b))
		if err != nil {
			if errors.Is(err, recomendation.ErrBudgetOutOfBounds) {
				httpError(ctx, w, http.StatusBadRequest, err)
			} else {
				httpError(ctx, w, http.StatusInternalServerError, err)
			}
			return
		}

		// Marshal response
		res, err := json.Marshal(rec)
		if err != nil {
			httpError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusOK)).Inc()
		_, _ = w.Write(res)
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		metrics.WritePrometheus(w, true)
	})

	return mux
}

func httpError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	slogctx.From(ctx).Error("request error", err)
	w.WriteHeader(status)
	metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, status)).Inc()
}

func recordDuration(startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	requestDuration.Update(duration)
}
