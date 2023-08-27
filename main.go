package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/izharishaksa/http-worker/config"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var conf *config.Config
var kafkaWriter *kafka.Writer
var kafkaReader *kafka.Reader

var (
	RequestIdKey = struct{}{}
)

func main() {
	var err error

	conf, err = config.NewConfig(time.Now)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-sigchan
		cancel()
	}()

	initKafkaReader(ctx)
	initKafkaWriter()
	consumeMessage(ctx)
}

func initKafkaWriter() {
	kafkaWriter = &kafka.Writer{
		Addr:                   kafka.TCP(conf.BrokerHost() + ":" + strconv.Itoa(conf.BrokerPort())),
		Topic:                  conf.DlqTopicName(),
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
}

func initKafkaReader(ctx context.Context) {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:     []string{conf.BrokerHost() + ":" + strconv.Itoa(conf.BrokerPort())},
		Topic:       conf.TopicName(),
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
		MaxAttempts: conf.MaxRetryAttempt(),
	}

	if conf.GroupId() != "" {
		kafkaConfig.GroupID = conf.GroupId()
	}

	kafkaReader = kafka.NewReader(kafkaConfig)
	if conf.GroupId() == "" {
		err := kafkaReader.SetOffsetAt(ctx, conf.OffsetTime())
		if err != nil {
			log.Fatalf("Failed to set offset: %v", err)
		}
	}
}

func consumeMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("context done, exiting")
			err := kafkaReader.Close()
			if err != nil {
				log.Fatalf("error closing kafka reader: %v", err)
			}
			return
		default:
			msg, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("error reading message: %v\n", err)
				continue
			}

			ctx = context.WithValue(ctx, RequestIdKey, getRequestIdFrom(msg.Headers))
			_ = processMessageWithRetry(ctx, msg)
		}
	}
}

func getRequestIdFrom(headers []kafka.Header) string {
	for _, header := range headers {
		if strings.ToLower(header.Key) == "x-request-id" {
			return string(header.Value)
		}
	}

	return uuid.New().String()
}

func processMessageWithRetry(ctx context.Context, msg kafka.Message) error {
	attempt := 0
	for attempt < conf.MaxRetryAttempt() {
		err := processMessage(ctx, string(msg.Value))
		if err != nil {
			log.Printf("error processing message: %v\n", err)
			attempt++
			time.Sleep(backoffDelay(attempt))
			continue
		}

		return nil
	}

	return publishToRetryTopic(ctx, msg)
}

func backoffDelay(attempt int) time.Duration {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	delay := baseDelay * (1 << uint(attempt))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

func processMessage(ctx context.Context, msg string) error {
	request, err := http.NewRequest(conf.InvokingMethod(), conf.InvokingUrl(), strings.NewReader(msg))
	if err != nil {
		log.Printf("error creating request: %v\n", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Request-Id", ctx.Value(RequestIdKey).(string))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("error sending request: %v\n", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("response received: %v\n", response.StatusCode)
		return err
	}

	log.Printf("response received: %v\n", response.StatusCode)

	return nil
}

func publishToRetryTopic(ctx context.Context, msg kafka.Message) error {
	msg.Topic = conf.DlqTopicName()
	err := kafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		log.Println("error publishing message to dlq topic", err.Error())
		return err
	}

	log.Println("message published to dlq topic")

	return nil
}
