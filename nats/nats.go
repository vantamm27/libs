package nats

import (
	"errors"
	"time"

	gonats "github.com/nats-io/go-nats"
)

type NatsConfig struct {
	Url   string `json:"url"`
	NumWk int    `json:"numwk"`
	Topic string `json:"topic"`
}

type NatsClient struct {
	Client *gonats.Conn
}

func (c *NatsClient) Publish(subject string, data []byte) error {
	if !c.Client.IsConnected() {
		return errors.New("status is disconnect")
	}
	return c.Publish(subject, data)

}

func NewClient(url string) (NatsClient, error) {

	client := NatsClient{}
	var err error

	client.Client, err = gonats.Connect(url, func(o *gonats.Options) error {
		o.Name = "Nats client"
		o.AllowReconnect = true
		o.MaxReconnect = 10
		o.ReconnectWait = 2 * time.Second
		return nil
	})

	return client, err

}
