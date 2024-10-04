package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go/jetstream"
	yopaNats "github.com/senorihl/yopa/pkg/nats"
	"github.com/senorihl/yopa/pkg/pixel"
	"os"
	"strings"
	"time"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	conf := yopaNats.Config{
		Name:  "Pixel receiver",
		Url:   os.Getenv("NATS_URL"),
		Pixel: struct{ Channel string }{Channel: os.Getenv("NATS_PIXEL_TOPIC")},
	}

	_, _, stream, _ := yopaNats.Setup(conf)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      "CONSUMER_" + conf.Pixel.Channel,
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	_, _ = c.Consume(func(msg jetstream.Msg) {
		if err := msg.Ack(); err != nil {
			log.Warn("Cannot acknowledge message: ", err)
		}

		parts := strings.Split(string(msg.Data()), "///")

		if event, err := pixel.UnparseQuery([]byte(parts[1])); err == nil {
			log.Info("Built event: ", event)
		} else {
			log.Info("Cannot build event: ", err)
		}
	})

	select {}

}
