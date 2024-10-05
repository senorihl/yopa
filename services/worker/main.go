package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Warn("Cannot connect to database: ", err)
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

			more, _ := json.Marshal(event.Event.More)
			row := db.QueryRow(`INSERT INTO yopa.public.pixel(
	  site_id, ts, visitor, name, 
	  page, page_chapter1, page_chapter2, page_chapter3, 
	  action, action_type, action_chapter1, action_chapter2, action_chapter3, 
	  custom_properties
  ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING (id);`,
				event.Site,
				event.Event.Globals.Timestamp,
				event.Event.Globals.Visitor,
				event.Event.Name,
				event.Event.Page.Name,
				event.Event.Page.Chapter1,
				event.Event.Page.Chapter2,
				event.Event.Page.Chapter3,
				event.Event.Action.Name,
				event.Event.Action.Type,
				event.Event.Action.Chapter1,
				event.Event.Action.Chapter2,
				event.Event.Action.Chapter3,
				more,
			)

			var id int64
			if err := row.Scan(&id); err != nil {
				log.Warn("Cannot save event due to ", err)
			} else {
				log.Infof("Saved event(%d): %s", id, event)
			}
		} else {
			log.Info("Cannot build event: ", err)
		}
	})

	select {}

}
