package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/senorihl/yopa/pkg/nats"
	"github.com/senorihl/yopa/services/pixel/server"
	"os"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	conf := nats.Config{
		Name:  "Pixel receiver",
		Url:   os.Getenv("NATS_URL"),
		Pixel: struct{ Channel string }{Channel: os.Getenv("NATS_PIXEL_TOPIC")},
	}

	nc, natsError := nats.Setup(conf)

	if natsError != nil {
		log.Warn("Cannot connect to NATS", natsError)
	}

	fiberApp := server.Setup(func(query string, remoteAddr string) {
		if natsError == nil {
			payload := []byte(fmt.Sprintf("%s///?%s", remoteAddr, query))
			err := nc.Publish(conf.Pixel.Channel, payload)

			if err == nil {
				log.Debug("Sent to NATS: ", string(payload))
			} else {
				log.Debug("Failed sending to NATS: ", string(payload), err)
			}
		} else {
			log.Warn("Didn't sent to NATS: ", natsError)
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
