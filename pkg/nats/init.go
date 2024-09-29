package nats

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	Name  string
	Url   string
	Pixel struct {
		Channel string
	}
}

func Setup(conf Config) (*nats.Conn, error) {

	opts := []nats.Option{nats.Name(conf.Name)}
	opts = setupConnOptions(opts)

	nc, err := nats.Connect(conf.Url, opts...)

	return nc, err
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
