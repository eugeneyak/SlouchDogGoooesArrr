// package bridge

// import (
// 	"github.com/nats-io/nats.go"
// 	"github.com/nats-io/nats.go/jetstream"
// )

// type Bridge struct {
// 	nc *nats.Conn
// }

// func Connect(url string) (*Bridge, error) {
// 	nc, err := nats.Connect(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	js, err := jetstream.New(nc)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = js.CreateStream(jetstream.StreamConfig{
// 		Name:      "tdlib.updates",
// 		Subjects:  []string{"tdlib.update"},
// 		MaxMsgs:   10000,
// 		Retention: jetstream.InterestPolicy,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	bridge := Bridge{
// 		nc: nc,
// 	}

// 	return &bridge, nil
// }

// func (bridge *Bridge) Drain() {
// 	bridge.nc.Drain()
// }

// func (bridge *Bridge) Emit(message Message) {

// }
