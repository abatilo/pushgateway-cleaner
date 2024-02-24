package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("PGC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Set up logging
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	log := slog.New(handler)

	pflag.BoolP(FlagVerbose, "v", FlagVerboseDefault, "Set logging verbosity level to debug")
	pflag.Duration(
		FlagTTL,
		FlagTTLDefault,
		"How old metrics are allowed to be before being deleted",
	)
	pflag.Duration(FlagSyncPeriod, FlagSyncPeriodDefault, "How often to check for old metrics")
	pflag.String(FlagPushgatewayURL, FlagPushgatewayURLDefault, "Pushgateway URL")
	pflag.BoolP(FlagHelp, "h", FlagHelpDefault, "Show help")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool(FlagHelp) {
		pflag.Usage()
		return
	}

	verbose := viper.GetBool(FlagVerbose)
	ttl := viper.GetDuration(FlagTTL)
	syncPeriod := viper.GetDuration(FlagSyncPeriod)
	pushgatewayURL := viper.GetString(FlagPushgatewayURL)

	if verbose {
		level.Set(slog.LevelDebug)
	}

	ticker := time.NewTicker(syncPeriod)

	log.Debug("Startup config",
		FlagVerbose, verbose,
		FlagTTL, ttl,
		FlagSyncPeriod, syncPeriod,
		FlagPushgatewayURL, pushgatewayURL,
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
		pushTimeMetric, err := fetchPushTimeMetric(httpClient, pushgatewayURL)
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
			metricLog.Debug("Found metric")

			if time.Since(pushTime) > ttl {
				metricLog.Debug(
					"Deleting metric since it's older than the TTL",
				)
				pushClient := push.New(pushgatewayURL, jobLabel)
				err := pushClient.
					Grouping("instance", instanceLabel).
					Client(httpClient).
					Delete()
				if err != nil {
					metricLog.Error("Error deleting metric:", err)
				}
			} else {
				metricLog.With(
					"time_until_delete", ttl-time.Since(pushTime),
				).Debug("Not deleting metric")
			}
		}
		log.Debug("Done fetching metrics")
	}
}
