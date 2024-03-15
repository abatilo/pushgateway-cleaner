package main

import (
	"fmt"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func fetchPushTimeMetric(
	httpClient HTTPClient,
	pushgatewayURL string,
) (*dto.MetricFamily, error) {
	resp, err := httpClient.Get(fmt.Sprintf("%s/metrics", pushgatewayURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser := new(expfmt.TextParser)
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	pushTimeMetrics, found := metricFamilies["push_time_seconds"]
	if !found {
		return &dto.MetricFamily{}, nil
	}

	return pushTimeMetrics, nil
}

func extractMetadata(
	metric *dto.Metric,
) (job string, instance string, pushTime time.Time) {
	for _, label := range metric.Label {
		if label.GetName() == "job" {
			job = label.GetValue()
		}
		if label.GetName() == "instance" {
			instance = label.GetValue()
		}
	}
	pushTime = time.Unix(int64(metric.GetGauge().GetValue()), 0)
	return
}
