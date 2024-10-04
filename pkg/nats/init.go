package nats

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

type Config struct {
	Name  string
	Url   string
	Pixel struct {
		Channel string
	}
}

func Setup(conf Config) (*nats.Conn, jetstream.JetStream, jetstream.Stream, error) {

	opts := []nats.Option{nats.Name(conf.Name)}
	opts = setupConnOptions(opts)

	nc, err := nats.Connect(conf.Url, opts...)

	if err != nil {
		return nil, nil, nil, err
	}

	js, jsError := jetstream.New(nc)

	if jsError != nil {
		return nc, js, nil, jsError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, streamErr := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     conf.Pixel.Channel,
		Subjects: []string{conf.Pixel.Channel + ".*"},
	})

	if streamErr != nil {
		log.Warn("Cannot get/create Stream in JetStream", streamErr)
	}

	return nc, js, stream, err
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Errorf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Errorf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))

	return opts
}
