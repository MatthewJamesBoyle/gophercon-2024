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
		defer func() {
			recordDuration(startTime)
		}()

		// Extract budget from query parameters and parse it
		budget := r.URL.Query().Get("budget")
		b, err := strconv.ParseInt(budget, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusBadRequest))
			return
		}

		// Get recommendation
		rec, err := svc.Get(ctx, int(b))
		if err != nil {
			switch {
			case errors.Is(err, recomendation.ErrBudgetOutOfBounds):
				w.WriteHeader(http.StatusBadRequest)
				metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusBadRequest))
			default:
				slogctx.From(ctx).Error("unhandled error", err)
				w.WriteHeader(http.StatusInternalServerError)
				metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusInternalServerError))
			}
			return
		}

		// Marshal response
		res, err := json.Marshal(rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusOK)
		metrics.GetOrCreateCounter(fmt.Sprintf(requestsTotal, http.StatusOK))
		_, _ = w.Write(res)
	})

	return mux
}

func recordDuration(startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	requestDuration.Update(duration)
}
