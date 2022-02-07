package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

func HttpRouter(conf Config, db *sql.DB) http.Handler {
	router := httptrace.New()
	router.Handler("GET", "/io-bound", &Handler{
		SQL: time.Duration(0.9 * float64(conf.Latency)),
		CPU: time.Duration(0.1 * float64(conf.Latency)),
		DB:  db,
	})
	router.Handler("GET", "/cpu-bound", &Handler{
		SQL: time.Duration(0.1 * float64(conf.Latency)),
		CPU: time.Duration(0.9 * float64(conf.Latency)),
		DB:  db,
	})
	return router
}

// Handler simulates an http request that spends the given amount of time
// in a SQL and On-CPU.
type Handler struct {
	SQL time.Duration
	CPU time.Duration
	DB  *sql.DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := h.DB.ExecContext(ctx, `SELECT pg_sleep($1)`, h.SQL.Seconds())
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusInternalServerError)
		return
	}

	cpuHog(h.CPU)
	fmt.Fprintf(w, "sql=%s cpu=%s\n", h.SQL, h.CPU)
}

func cpuHog(d time.Duration) {
	done := time.After(d)
	for {
		select {
		case <-done:
			return
		default:
			var m interface{}
			json.Unmarshal([]byte(`{"foo": [1, true, "bar"]}`), &m)
		}
	}
}
