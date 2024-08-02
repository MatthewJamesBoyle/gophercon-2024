package transporthttp

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"github.com/uhthomas/slogctx"
	"net/http"
	"strconv"
)

func NewMux(ctx context.Context, svc *recomendation.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /trip/{budget}", func(w http.ResponseWriter, r *http.Request) {
		budget := r.PathValue("budget")
		b, err := strconv.ParseInt(budget, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rec, err := svc.Get(ctx, int(b))
		if err != nil {
			switch {
			case errors.Is(err, recomendation.ErrBudgetOutOfBounds):
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				slogctx.From(ctx).Error("svc_get", "unhandled error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		res, err := json.Marshal(rec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(res)
	})

	return mux
}
