package bridge

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const Subject = "tdlib.update"

type Bridge struct {
	conn          *nats.Conn
	updatesStream jetstream.Stream
	methodsStream jetstream.Stream
}

func Connect(ctx context.Context, url string) (*Bridge, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	updatesStream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      "TDLibUpdates",
		Subjects:  []string{Subject},
		MaxMsgs:   10000,
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return nil, err
	}

	methodsStream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      "TDLibMethods",
		Subjects:  []string{"tdlib.method"},
		MaxMsgs:   10000,
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return nil, err
	}

	bridge := Bridge{
		conn:          conn,
		updatesStream: updatesStream,
		methodsStream: methodsStream,
	}

	return &bridge, nil
}

func (bridge *Bridge) Consume(ctx context.Context, handler func(msg jetstream.Msg)) error {
	consumer, err := bridge.methodsStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      "TDLibMethodsConsumer",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	consCtx, err := consumer.Consume(handler)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		consCtx.Stop()
		log.Println("STOP NATS CONS")
	}()

	return nil
}

func (bridge *Bridge) Drain() {
	bridge.conn.Drain()
}

func (bridge *Bridge) EmitUpdate(typ string, json []byte) {
	msg := nats.NewMsg(Subject)

	msg.Data = json
	msg.Header.Add("type", typ)

	err := bridge.conn.PublishMsg(msg)
	if err != nil {
		log.Println(err)
	}
}
