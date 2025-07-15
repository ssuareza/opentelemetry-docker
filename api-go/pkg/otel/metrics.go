package otel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

type metrics struct {
	activeRequests metric.Int64UpDownCounter
	httpRequests   metric.Int64Counter
	memoryUsage    metric.Int64ObservableGauge // system_memory_heap_bytes
}

// newMetrics initialize custom metrics
func newMetrics(meter metric.Meter) (*metrics, error) {
	var m metrics
	var err error

	// memory
	// system_memory_heap_bytes
	m.memoryUsage, err = meter.Int64ObservableGauge(
		"system.memory.heap",
		metric.WithDescription(
			"Memory usage of the allocated heap objects.",
		),
		metric.WithUnit("By"),
		metric.WithInt64Callback(
			func(ctx context.Context, o metric.Int64Observer) error {
				memoryUsage := getMemoryUsage()
				o.Observe(int64(memoryUsage))
				return nil
			},
		),
	)
	if err != nil {
		return nil, err
	}

	// http request counter

	m.httpRequests, err = meter.Int64Counter(
		"http.server.requests",
		metric.WithDescription("Total number of HTTP requests received."),
		metric.WithUnit("{requests}"),
	)
	if err != nil {
		return nil, err
	}

	// http active request counter
	// 	m.activeRequests, err = meter.Int64UpDownCounter(
	// 		"http.server.active_requests",
	// 		metric.WithDescription("Number of in-flight requests."),
	// 		metric.WithUnit("{requests}"),
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	return &m, nil
}

// getMemoryUsage returns the current memory usage of the process.
func getMemoryUsage() uint64 {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	currentMemoryUsage := memStats.HeapAlloc

	return currentMemoryUsage
}
