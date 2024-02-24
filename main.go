package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	log := slog.New(handler)

	verbose := flag.Bool(FlagVerbose, FlagVerboseDefault, "Set logging verbosity level to debug")
	ttl := flag.Duration(
		FlagTTL,
		FlagTTLDefault,
		"How old metrics are allowed to be before being deleted",
	)
	syncPeriod := flag.Duration(
		FlagSyncPeriod,
		FlagSyncPeriodDefault,
		"How often to check for old metrics",
	)
	pushgatewayURL := flag.String(FlagPushgatewayURL, FlagPushgatewayURLDefault, "Pushgateway URL")
	flag.Parse()

	if *verbose {
		level.Set(slog.LevelDebug)
	}

	ticker := time.NewTicker(*syncPeriod)

	log.Debug("Startup config",
		FlagVerbose, *verbose,
		FlagTTL, *ttl,
		FlagSyncPeriod, *syncPeriod,
		FlagPushgatewayURL, *pushgatewayURL,
	)

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 3 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	for range ticker.C {
		log.Debug("Fetching metrics")
		pushTimeMetric, err := fetchPushTimeMetric(httpClient, *pushgatewayURL)
		if err != nil {
			log.Error("Error fetching or parsing metrics:", err)
			return
		}

		for _, metric := range pushTimeMetric.GetMetric() {
			jobLabel, instanceLabel, pushTime := extractMetadata(metric)
			metricLog := log.With(
				"instance", instanceLabel,
				"job", jobLabel,
				"push_time", pushTime,
			)
			metricLog.Debug("Found metric group")

			if time.Since(pushTime) > *ttl {
				metricLog.Debug(
					"Deleting metric group since it's older than the TTL",
				)
				url, _ := url.Parse(fmt.Sprintf(
					"%s/metrics/job/%s/instance/%s",
					*pushgatewayURL,
					jobLabel,
					instanceLabel,
				))
				_, err := httpClient.Do(&http.Request{
					Method: http.MethodDelete,
					URL:    url,
				})
				if err != nil {
					metricLog.Error("Error deleting metric group:", err)
				}
			} else {
				metricLog.With(
					"time_until_delete", *ttl-time.Since(pushTime),
				).Debug("Not deleting metric group")
			}
		}
		log.Debug("Done fetching metrics")
	}
}
