package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type monitor struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	rateLimit      int
}

type monitorOptions struct {
	PollIntervalSeconds   int
	ReportIntervalSeconds int
	RateLimit             int
}

func newMonitor(opts monitorOptions) *monitor {
	return &monitor{
		pollInterval:   time.Duration(opts.PollIntervalSeconds) * time.Second,
		reportInterval: time.Duration(opts.ReportIntervalSeconds) * time.Second,
		rateLimit:      opts.RateLimit,
	}
}

func (m *monitor) Start(ctx context.Context, client *client) {
	tasks := m.startMetricsCollector(ctx)
	wg := m.startMetricsSender(ctx, client, tasks)
	log.Debug().Msg("waiting for workers to finish")
	wg.Wait()
}

func (m *monitor) startMetricsCollector(ctx context.Context) chan metrics.Metric {
	out := make(chan metrics.Metric)

	go func() {
		defer func() {
			log.Debug().Msg("close metrics collector output channel")
			close(out)
		}()

		ticker := time.NewTicker(m.pollInterval)
		defer ticker.Stop()

		var counter int64
		mm := make(map[string]metrics.Metric)

		for {
			select {
			case <-ticker.C:
				log.Debug().Msg("collecting metrics")
				counter++
				mm["PollCount"] = newCounter("PollCount", counter)
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				mm["Alloc"] = newGauge("Alloc", float64(memStats.Alloc))
				mm["BuckHashSys"] = newGauge("BuckHashSys", float64(memStats.BuckHashSys))
				mm["Frees"] = newGauge("Frees", float64(memStats.Frees))
				mm["GCCPUFraction"] = newGauge("GCCPUFraction", memStats.GCCPUFraction)
				mm["GCSys"] = newGauge("GCSys", float64(memStats.GCSys))
				mm["HeapAlloc"] = newGauge("HeapAlloc", float64(memStats.HeapAlloc))
				mm["HeapIdle"] = newGauge("HeapIdle", float64(memStats.HeapIdle))
				mm["HeapInuse"] = newGauge("HeapInuse", float64(memStats.HeapInuse))
				mm["HeapObjects"] = newGauge("HeapObjects", float64(memStats.HeapObjects))
				mm["HeapReleased"] = newGauge("HeapReleased", float64(memStats.HeapReleased))
				mm["HeapSys"] = newGauge("HeapSys", float64(memStats.HeapSys))
				mm["LastGC"] = newGauge("LastGC", float64(memStats.LastGC))
				mm["Lookups"] = newGauge("Lookups", float64(memStats.Lookups))
				mm["MCacheInuse"] = newGauge("MCacheInuse", float64(memStats.MCacheInuse))
				mm["MCacheSys"] = newGauge("MCacheSys", float64(memStats.MCacheSys))
				mm["MSpanInuse"] = newGauge("MSpanInuse", float64(memStats.MSpanInuse))
				mm["MSpanSys"] = newGauge("MSpanSys", float64(memStats.MSpanSys))
				mm["Mallocs"] = newGauge("Mallocs", float64(memStats.Mallocs))
				mm["NextGC"] = newGauge("NextGC", float64(memStats.NextGC))
				mm["NumForcedGC"] = newGauge("NumForcedGC", float64(memStats.NumForcedGC))
				mm["NumGC"] = newGauge("NumGC", float64(memStats.NumGC))
				mm["OtherSys"] = newGauge("OtherSys", float64(memStats.OtherSys))
				mm["PauseTotalNs"] = newGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
				mm["StackInuse"] = newGauge("StackInuse", float64(memStats.StackInuse))
				mm["StackSys"] = newGauge("StackSys", float64(memStats.StackSys))
				mm["Sys"] = newGauge("Sys", float64(memStats.Sys))
				mm["TotalAlloc"] = newGauge("TotalAlloc", float64(memStats.TotalAlloc))
				mm["RandomValue"] = newGauge("RandomValue", rand.Float64())
				virtualMemStats, _ := mem.VirtualMemory()
				mm["TotalMemory"] = newGauge("TotalMemory", float64(virtualMemStats.Total))
				mm["FreeMemory"] = newGauge("FreeMemory", float64(virtualMemStats.Free))
				cpuPercentages, _ := cpu.Percent(time.Second, true)
				for i, percent := range cpuPercentages {
					key := fmt.Sprintf("CPUutilization%d", i)
					mm[key] = newGauge(key, percent)
				}
				for _, m := range mm {
					out <- m
				}
			case <-ctx.Done():
				log.Debug().Msg("context done, stopping metrics collector")
				return
			}
		}
	}()
	return out
}

func (m *monitor) startMetricsSender(ctx context.Context, client *client, in chan metrics.Metric) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	for n := range m.rateLimit {
		wg.Add(1)
		go func(n int) {
			log.Debug().Msgf("starting metric sender %d", n)
			defer wg.Done()
			for {
				select {
				case metric := <-in:
					if err := client.SendOne(ctx, metric); err != nil {
						log.Error().Err(err).Msg("failed to send metric")
					}
				case <-ctx.Done():
					log.Debug().Msg("context done, stopping metric sender")
					return
				}
			}
		}(n)
	}
	return wg
}

func newCounter(id string, n int64) metrics.Metric {
	value := n
	return metrics.Metric{ID: id, Type: metrics.Counter, Counter: &value}
}

func newGauge(id string, n float64) metrics.Metric {
	value := n
	return metrics.Metric{ID: id, Type: metrics.Gauge, Gauge: &value}
}
