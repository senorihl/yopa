package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	yopaNats "github.com/senorihl/yopa/pkg/nats"
	"github.com/senorihl/yopa/pkg/pixel"
	"os"
	"runtime"
	"strings"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	conf := yopaNats.Config{
		Name:  "Pixel receiver",
		Url:   os.Getenv("NATS_URL"),
		Pixel: struct{ Channel string }{Channel: os.Getenv("NATS_PIXEL_TOPIC")},
	}

	nc, _ := yopaNats.Setup(conf)

	nc.Subscribe(conf.Pixel.Channel, func(m *nats.Msg) {
		parts := strings.Split(string(m.Data), "///")
		log.Info("Received a message: ", parts)
		if event, err := pixel.UnparseQuery([]byte(parts[1])); err == nil {
			log.Info("Built event: ", event)
		} else {
			log.Info("Cannot build event: ", err)
		}
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Infof("Listening on [%s]", conf.Pixel.Channel)

	runtime.Goexit()
}
