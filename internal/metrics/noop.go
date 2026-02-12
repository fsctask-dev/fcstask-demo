package metrics

type NoopMetricsClient struct{}

func (NoopMetricsClient) Inc(name string, tags map[string]string)                  {}
func (NoopMetricsClient) Add(name string, v float64, tags map[string]string)       {}
func (NoopMetricsClient) Gauge(name string, v float64, tags map[string]string)     {}
func (NoopMetricsClient) Histogram(name string, v float64, tags map[string]string) {}
func (NoopMetricsClient) Close() error                                             { return nil }
