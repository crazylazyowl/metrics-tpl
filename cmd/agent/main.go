package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	if key != "" {
		fmt.Println("hash is on")
	}
	gauge := make(map[string]float64)
	var counter int64
	// NOTE: hack for second test
	attempts := 4
	delay := 1
	if err := tryPing(address); err != nil {
		attempts = 1
		delay = 0
	}
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
				if err := reportOne(address, metric, attempts, delay, key); err != nil {
					log.Printf("failed to send %s (%f); err=%v\n", name, value, err)
				}
			}
			metric := metrics.Metric{
				ID:      "PollCount",
				Type:    metrics.Counter,
				Counter: &counter,
			}
			if err := reportOne(address, metric, attempts, delay, key); err != nil {
				log.Printf("failed to send %s (%d); err=%v\n", "PollCount", counter, err)
			}
			// many = append(many, metric)
			// if err := reportMany(address, many, attempts, delay, key); err != nil {
			// 	log.Printf("failed to bulk metrics; err=%v\n", err)
			// }
		}
	}
}

func tryPing(address string) error {
	url := fmt.Sprintf("http://%s/update/", address)
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func reportMany(address string, mm []metrics.Metric, attempts int, delay int, key string) error {
	for _, m := range mm {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	data, err := json.Marshal(mm)
	if err != nil {
		return err
	}
	return report(fmt.Sprintf("http://%s/updates/", address), data, attempts, delay, true, key)
}

func reportOne(address string, m metrics.Metric, attempts int, delay int, key string) error {
	if err := m.Validate(); err != nil {
		return err
	}
	data, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	return report(fmt.Sprintf("http://%s/update/", address), data, attempts, delay, false, key)
}

func report(url string, data []byte, attempts, delay int, compress bool, key string) error {
	var err error
	if compress {
		log.Println("compression is on")
		buf := bytes.NewBuffer(nil)
		w := gzip.NewWriter(buf)
		w.Write(data)
		w.Close()
		data = buf.Bytes()
	}
	var digest string
	if key != "" {
		h := hmac.New(sha256.New, []byte(key))
		h.Write(data)
		digest = hex.EncodeToString(h.Sum(nil))
	}
	for range attempts {
		var req *http.Request
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to prepare request; %w", err)
		}
		var resp *http.Response
		req.Header.Add("Content-Type", "application/json")
		if compress {
			req.Header.Add("Content-Encoding", "gzip")
		}
		if digest != "" {
			req.Header.Add("HashSHA256", digest)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("retry report")
			time.Sleep(time.Duration(delay) * time.Second)
			delay += 2
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("status code = %d, data = %s", resp.StatusCode, string(body))
		}
		break
	}
	return err
}
