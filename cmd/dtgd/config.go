package main

import (
	"flag"
	"time"
)

func ConfigFromFlags() Config {
	conf := Config{Version: version}
	flag.StringVar(&conf.Addr, "addr", "localhost:9191", "HTTP addr to listen on.")
	flag.StringVar(&conf.Env, "env", "dev", "The env tag for Datadog.")
	flag.StringVar(&conf.Service, "service", "dtgd", "The service tag for Datadog.")
	flag.DurationVar(&conf.Latency, "latency", 100*time.Millisecond, "The request response time to simulate.")
	flag.StringVar(&conf.DB, "db", "postgres://dtgd:dtgd-secret@localhost:5432/", "The dsn for the postgresql connection.")
	flag.Parse()
	return conf
}

// Config controls the behavior of the demo application.
type Config struct {
	Addr    string
	Service string
	Env     string
	Version string
	DB      string
	Latency time.Duration
}
