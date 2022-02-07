package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var version = "0.1"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	conf := ConfigFromFlags()
	log.Printf("Starting %s version %s", conf.Service, conf.Version)

	sqltrace.Register("pgx", stdlib.GetDefaultDriver())

	stopTracer := startTracer(conf)
	defer stopTracer()

	stopProfiler, err := startProfiler(conf)
	if err != nil {
		return err
	}
	defer stopProfiler()

	db, err := openDB(conf)
	if err != nil {
		return err
	}
	return serveHttp(conf, db)
}

func startTracer(conf Config) func() {
	tracer.Start(
		tracer.WithEnv(conf.Env),
		tracer.WithService(conf.Service),
		tracer.WithServiceVersion(version),
		tracer.WithGlobalTag("go_version", runtime.Version()),
		tracer.WithProfilerCodeHotspots(true),
		tracer.WithProfilerEndpoints(true),
	)
	return tracer.Stop
}

func startProfiler(conf Config) (func(), error) {
	profilerOptions := []profiler.Option{
		// Important: CPUDuration should match Period to achieve 100% Code Hotspots
		// coverage. Default is 25% right now, but this might change in the next
		// dd-trace-go release.
		profiler.CPUDuration(60 * time.Second),
		profiler.WithPeriod(60 * time.Second),

		profiler.WithService(conf.Service),
		profiler.WithEnv(conf.Env),
		profiler.WithVersion(conf.Version),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.BlockProfile,
			profiler.MutexProfile,
			profiler.GoroutineProfile,
		),
		profiler.WithTags("go_version:" + runtime.Version()),
	}

	return profiler.Stop, profiler.Start(profilerOptions...)
}

func openDB(conf Config) (*sql.DB, error) {
	db, err := sqltrace.Open("pgx", conf.DB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func serveHttp(conf Config, db *sql.DB) error {
	log.Printf("Serving on http://%s/", conf.Addr)
	return http.ListenAndServe(conf.Addr, HttpRouter(conf, db))
}
