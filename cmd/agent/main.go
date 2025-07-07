package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

const (
// baseURL    = "http://localhost:8080/update"
// gaugeURL   = "http://%s/update/gauge/%s/%f"
// counterURL = "http://%s/update/counter/%s/%d"
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
	var counter int64
	url := fmt.Sprintf("http://%s/update/", conf.address)
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
			// many := make([]metrics.Metric, 0, len(gauge)+1)
			for key, value := range gauge {
				metric := metrics.Metric{
					ID:    key,
					Type:  metrics.Gauge,
					Gauge: &value,
				}
				// many = append(many, metric)
				if err := report(url, &metric); err != nil {
					log.Printf("failed to send %s (%f); err=%v\n", key, value, err)
				}
			}
			metric := metrics.Metric{
				ID:      "PollCount",
				Type:    metrics.Counter,
				Counter: &counter,
			}
			if err := report(url, &metric); err != nil {
				log.Printf("failed to send %s (%d); err=%v\n", "PollCount", counter, err)
			}
			// many = append(many, metric)
			// if err := reportBulk(url, many); err != nil {
			// 	log.Printf("failed to bulk metrics; err=%v\n", err)
			// }
		}
	}
}

func report(url string, metric *metrics.Metric) error {
	// if err := metric.Validate(); err != nil {
	// 	return err
	// }

	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	// buf := bytes.NewBuffer(nil)

	// w := gzip.NewWriter(buf)
	// w.Write(data)
	// w.Close()

	// req, _ := http.NewRequest(http.MethodPost, url, buf)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code = %d, data = %s", resp.StatusCode, string(body))
	}

	return nil
}

// func reportBulk(url string, many []metrics.Metric) error {
// 	// buf := bytes.NewBuffer(nil)
// 	// w := gzip.NewWriter(buf)
// 	// if err := json.NewEncoder(w).Encode(many); err != nil {
// 	// 	return err
// 	// }
// 	// w.Close()

// 	data, err := json.Marshal(many)
// 	if err != nil {
// 		return err
// 	}

// 	// log.Printf(string(data))

// 	// req, _ := http.NewRequest(http.MethodPost, url, buf)
// 	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 	req.Header.Add("Content-Type", "application/json")
// 	// req.Header.Add("Content-Encoding", "gzip")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	body, _ := io.ReadAll(resp.Body)
// 	resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return fmt.Errorf("status code = %d, data = %s", resp.StatusCode, string(body))
// 	}

// 	return nil

// }
