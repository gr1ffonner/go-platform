package metrics

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemMetrics holds system resource metrics
type SystemMetrics struct {
	MemoryUsage     prometheus.Gauge
	CPUUsage        prometheus.Gauge
	GoroutinesCount prometheus.Gauge
}

// NewSystemMetrics creates a new system metrics instance
func NewSystemMetrics(registry *prometheus.Registry) *SystemMetrics {
	return &SystemMetrics{
		MemoryUsage: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "go_memstats_heap_inuse_bytes",
				Help: "Go heap memory in use",
			},
		),

		CPUUsage: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "go_cpu_usage_percent",
				Help: "Go application CPU usage percentage",
			},
		),

		GoroutinesCount: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "go_goroutines",
				Help: "Number of goroutines",
			},
		),
	}
}

// StartCollection starts collecting system metrics in the background
func (s *SystemMetrics) StartCollection(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	slog.Info("Starting system metrics collection", "interval", "15s")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping system metrics collection")
			return
		case <-ticker.C:
			s.collectSystemMetrics()
		}
	}
}

// collectSystemMetrics collects current system metrics
func (s *SystemMetrics) collectSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Set memory usage (heap in use)
	s.MemoryUsage.Set(float64(memStats.HeapInuse))

	// Set goroutines count
	s.GoroutinesCount.Set(float64(runtime.NumGoroutine()))

	// Set CPU usage
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err == nil {
		cpuPercent, err := proc.CPUPercent()
		if err == nil {
			s.CPUUsage.Set(cpuPercent)
		}
	}
}
