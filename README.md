# dd-trace-go-demo

A simple application to show how to use [dd-trace-go](https://github.com/DataDog/dd-trace-go)'s tracer and profiler.

## Usage

To run this demo application, simply clone this repository and follow the steps outlined below:

```
# set the api key of your datadog account
export DD_API_KEY=...

# launch datadog agent and postgres
docker-compose up

# launch go application
go run ./cmd/dtgd/

# hit endpoints
curl localhost:9191/cpu-bound
curl localhost:9191/io-bound
```

Note: It might take a few seconds for traces to show up, and a little over 1 minute for the first profile to be uploaded.
