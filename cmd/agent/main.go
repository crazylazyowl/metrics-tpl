package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := newClient(clientOptions{
		BaseURL: conf.address,
		Secret:  conf.key,
	})
	client.Hack(ctx)

	monitor := newMonitor(monitorOptions{
		PollIntervalSeconds:   conf.pollInterval,
		ReportIntervalSeconds: conf.reportInterval,
		RateLimit:             conf.rateLimit,
	})

	if err := monitor.Start(ctx, client); err != nil {
		log.Fatalln(err)
	}
}
