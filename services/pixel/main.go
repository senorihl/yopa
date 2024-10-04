package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/senorihl/yopa/pkg/nats"
	"github.com/senorihl/yopa/services/pixel/server"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	conf := nats.Config{
		Name:  "Pixel receiver",
		Url:   os.Getenv("NATS_URL"),
		Pixel: struct{ Channel string }{Channel: os.Getenv("NATS_PIXEL_TOPIC")},
	}

	_, js, _, err := nats.Setup(conf)

	if err != nil {
		log.Warn("Cannot setup NATS with Jetstream: ", err)
	}

	fiberApp := server.Setup(func(query string, remoteAddr string) {
		if js != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			payload := []byte(fmt.Sprintf("%s///?%s", remoteAddr, query))
			_, err := js.Publish(ctx, conf.Pixel.Channel+".new", payload)

			if err == nil {
				log.Debug("Sent to JetStream: ", string(payload))
			} else {
				log.Debug("Failed sending to JetStream: ", string(payload), err)
			}
		}
	})

	fiberApp.Hooks().OnListen(func(data fiber.ListenData) error {
		log.Infof("Listening on [http://%s:%s]", data.Host, data.Port)
		return nil
	})

	fiberError := fiberApp.Listen(":80")
	if fiberError != nil {
		log.Errorf("Failed to start webserver: %s", fiberError)
		return
	}
}
