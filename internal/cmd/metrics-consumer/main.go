package metrics_consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	"fcstask/internal/config"
	kafkaMetric "fcstask/internal/kafka"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          cfg.Kafka.TopicMetrics,
		GroupID:        "metrics-consumer",
		MinBytes:       cfg.Kafka.MinBytes,
		MaxBytes:       cfg.Kafka.MaxBytes,
		MaxWait:        cfg.Kafka.ReadTimeout,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	defer reader.Close()

	log.Printf("[Metrics Consumer] Started. Topic: %s", cfg.Kafka.TopicMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	go func() {
		<-sigchan
		log.Println("[Metrics Consumer] Shutting down...")
		cancel()
	}()

	for ctx.Err() == nil {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Println("[Metrics Consumer] Error reading message:", err)
			continue
		}

		var metric kafkaMetric.Metric
		err = json.Unmarshal(msg.Value, &metric)
		if err != nil {
			log.Println("[Metrics Consumer] Error unmarshalling message:", err)
			continue
		}

		//пока какой-то лог, потом это всё в бд будет кидаться
		tags := ""
		if len(metric.Tags) > 0 {
			tagsBytes, _ := json.Marshal(metric.Tags)
			tags = string(tagsBytes)
		}

		log.Printf("[METRIC] time=%s name=%s value=%.2f tags=%s",
			metric.Timestamp.Format(time.RFC3339),
			metric.Name,
			metric.Value,
			tags,
		)
		}
	}

}
