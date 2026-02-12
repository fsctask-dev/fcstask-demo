# Пакет metrics-consumer
## cmd/metrics-consumer/
```text
cmd/metrics-consumer/
└── main.go  # Consumer метрик из Kafka
```
## Назначение
Читает метрики из Kafka топика и выводит в консоль (или сохраняет в БД).

```bash
go run cmd/metrics-consumer/main.go
```
Пример вывода:

```text
[14:25:31.123] http_requests_total = 1.00 {"path":"/hello","method":"GET","status":"200"}
[14:25:31.124] http_request_duration_seconds = 0.002 {"path":"/hello","method":"GET"}
[14:25:31.125] tasks_submitted_total = 1.00 {"course":"algorithms","assignment":"hw1"}
```
Конфигурация
```yaml
kafka:
brokers: ["localhost:9092"]
topic_metrics: "edu.metrics"
group_id: "metrics-consumer"  # consumer group
```
Параметры:

* brokers — список брокеров Kafka

* topic_metrics — топик с метриками

* group_id — идентификатор consumer group

## Расширение
Вывод лога в консоль будет заменён на запись в PostgreSQL

## Пример полного flow
```go
// 1. API сервис
producer := kafka.NewProducer(cfg)
client := metrics.NewKafkaMetricsClient(producer, cfg.TopicMetrics, ctx)

metricsMW := middleware.NewMetricsMiddleware(client)
inflightMW := middleware.NewInflightMetricsMiddleware(client)
business := metrics.NewBusinessMetrics(client)

mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
business.TaskSubmitted("algorithms", "hw1")
w.WriteHeader(http.StatusAccepted)
})

http.ListenAndServe(":8080", inflightMW.Handler(metricsMW.Handler(mux)))

// 2. Consumer (отдельный процесс)
consumer := kafka.NewReader(kafka.ReaderConfig{
Brokers: cfg.Brokers,
Topic:   cfg.TopicMetrics,
})

for {
msg, _ := consumer.ReadMessage(ctx)
var metric kafka.Metric
json.Unmarshal(msg.Value, &metric)
saveToDatabase(metric)  // или вывести в консоль
}
```
