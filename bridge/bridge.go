package bridge

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const Subject = "tdlib.update"

type Bridge struct {
	nc *nats.Conn
}

func Connect(ctx context.Context, url string) (*Bridge, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	config := jetstream.StreamConfig{
		Name:      "TDLibUpdates",
		Subjects:  []string{Subject},
		MaxMsgs:   10000,
		Retention: jetstream.InterestPolicy,
	}

	_, err = js.CreateStream(ctx, config)
	if err != nil {
		return nil, err
	}

	bridge := Bridge{
		nc: nc,
	}

	return &bridge, nil
}

func (bridge *Bridge) Drain() {
	bridge.nc.Drain()
}

func (bridge *Bridge) EmitUpdate(typ string, json []byte) {
	msg := nats.NewMsg(Subject)

	msg.Data = json
	msg.Header.Add("type", typ)

	err := bridge.nc.PublishMsg(msg)
	if err != nil {
		log.Println(err)
	}
}
