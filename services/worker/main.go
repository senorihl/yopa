package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	yopaNats "github.com/senorihl/yopa/pkg/nats"
	"os"
	"runtime"
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
		log.Info("Received a message: ", string(m.Data))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Infof("Listening on [%s]", conf.Pixel.Channel)

	runtime.Goexit()
}
