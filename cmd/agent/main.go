package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	monitor(ctx, conf.address, conf.pollInterval, conf.reportInterval, conf.key)
}

func monitor(ctx context.Context, address string, pollInterval, reportInterval int, key string) error {
	client := newClient(clientOptions{
		BaseURL: address,
		Secret:  key,
	})
	client.Hack(ctx)

	gauge := make(map[string]float64)
	var counter int64

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			log.Println("read mem stats")
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			gauge["Alloc"] = float64(ms.Alloc)
			gauge["BuckHashSys"] = float64(ms.BuckHashSys)
			gauge["Frees"] = float64(ms.Frees)
			gauge["GCCPUFraction"] = ms.GCCPUFraction
			gauge["GCSys"] = float64(ms.GCSys)
			gauge["HeapAlloc"] = float64(ms.HeapAlloc)
			gauge["HeapIdle"] = float64(ms.HeapIdle)
			gauge["HeapInuse"] = float64(ms.HeapInuse)
			gauge["HeapObjects"] = float64(ms.HeapObjects)
			gauge["HeapReleased"] = float64(ms.HeapReleased)
			gauge["HeapSys"] = float64(ms.HeapSys)
			gauge["LastGC"] = float64(ms.LastGC)
			gauge["Lookups"] = float64(ms.Lookups)
			gauge["MCacheInuse"] = float64(ms.MCacheInuse)
			gauge["MCacheSys"] = float64(ms.MCacheSys)
			gauge["MSpanInuse"] = float64(ms.MSpanInuse)
			gauge["MSpanSys"] = float64(ms.MSpanSys)
			gauge["Mallocs"] = float64(ms.Mallocs)
			gauge["NextGC"] = float64(ms.NextGC)
			gauge["NumForcedGC"] = float64(ms.NumForcedGC)
			gauge["NumGC"] = float64(ms.NumGC)
			gauge["OtherSys"] = float64(ms.OtherSys)
			gauge["PauseTotalNs"] = float64(ms.PauseTotalNs)
			gauge["StackInuse"] = float64(ms.StackInuse)
			gauge["StackSys"] = float64(ms.StackSys)
			gauge["Sys"] = float64(ms.Sys)
			gauge["TotalAlloc"] = float64(ms.TotalAlloc)
			gauge["RandomValue"] = rand.Float64()
			vm, _ := mem.VirtualMemory()
			gauge["TotalMemory"] = float64(vm.Total)
			gauge["FreeMemory"] = float64(vm.Free)
			percentages, _ := cpu.Percent(time.Second, true)
			for i, percent := range percentages {
				key := fmt.Sprintf("CPUutilization%d", i)
				gauge[key] = percent
			}
			counter++
		case <-reportTicker.C:
			log.Println("send metrics")
			// many := make([]metrics.Metric, 0, len(gauge)+1)
			for name, value := range gauge {
				metric := metrics.Metric{
					ID:    name,
					Type:  metrics.Gauge,
					Gauge: &value,
				}
				// many = append(many, metric)
				if err := client.SendOne(ctx, metric); err != nil {
					log.Printf("failed to send %s (%f); err=%v\n", name, value, err)
				}
			}
			metric := metrics.Metric{
				ID:      "PollCount",
				Type:    metrics.Counter,
				Counter: &counter,
			}
			if err := client.SendOne(ctx, metric); err != nil {
				log.Printf("failed to send %s (%d); err=%v\n", "PollCount", counter, err)
			}
			// many = append(many, metric)
			// if err := reportMany(address, many, attempts, delay, key); err != nil {
			// 	log.Printf("failed to bulk metrics; err=%v\n", err)
			// }
		}
	}
}
