package transporthttp

import (
	"encoding/json"
	"errors"
	"github.com/matthewjamesboyle/gophercon-2024/internal/recomendation"
	"net/http"
	"strconv"
)

func NewMux(svc *recomendation.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /trip/{budget}", func(w http.ResponseWriter, r *http.Request) {
		budget := r.PathValue("budget")
		b, err := strconv.ParseInt(budget, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rec, err := svc.Get(r.Context(), int(b))
		if err != nil {
			switch {
			case errors.Is(err, recomendation.ErrBudgetOutOfBounds):
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
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
