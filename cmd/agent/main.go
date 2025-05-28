package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	// baseURL    = "http://localhost:8080/update"
	gaugeURL   = "http://%s/update/gauge/%s/%f"
	counterURL = "http://%s/update/counter/%s/%d"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	_ = monitor(ctx, conf)
}

func monitor(ctx context.Context, conf *config) error {
	gauge := make(map[string]float64)
	var counter uint64

	pollTicker := time.NewTicker(time.Duration(conf.pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(conf.reportInterval) * time.Second)
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
			counter++
		case <-reportTicker.C:
			log.Println("send metrics")
			for key, value := range gauge {
				if err := report(fmt.Sprintf(gaugeURL, conf.address, key, value)); err != nil {
					log.Printf("failed to send %s (%f); err=%v\n", key, value, err)
				}
			}
			if err := report(fmt.Sprintf(counterURL, conf.address, "PollCount", counter)); err != nil {
				log.Printf("failed to send %s (%d); err=%v\n", "PollCount", counter, err)
			}
		}
	}
}

func report(url string) error {
	resp, err := http.Post(url, "text/plain", nil)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("status code = %d", resp.StatusCode)
		}
	}
	return err
}
