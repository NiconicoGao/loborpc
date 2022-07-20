package lobo

import (
	"context"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	file  *os.File
	kafka kafka.Writer
}

func NewKafkaWriter() *KafkaWriter {
	file, _ := os.Create("./test.log")
	k := kafka.Writer{
		Addr:         kafka.TCP([]string{"52.11.26.186:9092"}...),
		Topic:        "logger",
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: time.Second * 3,
	}

	return &KafkaWriter{
		file:  file,
		kafka: k,
	}
}

func (c *KafkaWriter) Write(p []byte) (n int, err error) {
	err = c.kafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: p,
		},
	)

	if err != nil {
		return 0, err
	}

	return c.file.Write(p)
}
