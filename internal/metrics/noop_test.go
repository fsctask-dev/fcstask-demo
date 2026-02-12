package metrics_test

import (
	"testing"

	"fcstask/internal/metrics"
)

func TestNoopMetricsClient(t *testing.T) {
	client := metrics.NoopMetricsClient{}

	// Проверяем, что ничего не падает
	client.Inc("test", nil)
	client.Add("test", 5, nil)
	client.Gauge("test", 1, nil)
	client.Histogram("test", 0.5, nil)

	if err := client.Close(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestNoopMetricsClient_WithTags(t *testing.T) {
	client := metrics.NoopMetricsClient{}

	// С тегами тоже не должно падать
	tags := map[string]string{
		"env":  "test",
		"app":  "fcstask",
		"path": "/api/v1/test",
	}

	client.Inc("test_counter", tags)
	client.Gauge("test_gauge", 100, tags)
	client.Histogram("test_histogram", 0.123, tags)
}
