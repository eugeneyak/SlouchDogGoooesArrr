package tdlib

/*
#cgo CFLAGS: -I/home/yak/td/tdlib/include
#cgo LDFLAGS: -L/home/yak/td/build -Wl,-rpath,/home/yak/td/build -ltdjson -ldl

#include <stdlib.h>
#include <td/telegram/td_json_client.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"log"
	"unsafe"
)

// Client wraps a TDLib JSON client instance.
type Client struct {
	ptr unsafe.Pointer
	log *log.Logger
}

// Creates a new TDLib JSON client instance.
func Init() *Client {
	log := log.New(log.Writer(), "[TDLib] ", log.LstdFlags|log.Lshortfile)
	ptr := C.td_json_client_create()

	log.Println("TDLib client initialized:", ptr)

	return &Client{
		ptr: unsafe.Pointer(ptr),
		log: log,
	}
}

// Send sends a JSON-formatted request to the TDLib client.
func (c *Client) Send(request string) {
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))

	C.td_json_client_send(c.ptr, cRequest)
}

// Execute synchronously executes a TDLib request.
// The returned string is valid until the next call to Receive or Execute.
func (c *Client) Execute(request string) string {
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))

	result := C.td_json_client_execute(c.ptr, cRequest)
	if result == nil {
		return ""
	}

	return C.GoString(result)
}

// Receive waits for a TDLib response or update for up to timeout seconds.
// It returns the JSON response string or an empty string on timeout.
func (c *Client) Receive(ctx context.Context) chan Update {
	channel := make(chan Update)
	go c.receive(ctx, channel)

	return channel
}

func (c *Client) receive(ctx context.Context, updates chan Update) {
	for {
		select {
		case <-ctx.Done():
			close(updates)

		default:
			raw := C.td_json_client_receive(c.ptr, C.double(10))
			if raw == nil {
				continue
			}

			str := C.GoString(raw)

			var update Update

			if err := json.Unmarshal([]byte(str), &update); err == nil {
				update.Payload = str
				updates <- update
			} else {
				c.log.Printf("Failed to unmarshal update: %s", err.Error())
			}
		}
	}
}

// Destroy releases the TDLib client instance.
func (c *Client) Destroy() {
	C.td_json_client_destroy(c.ptr)
	c.ptr = nil
}
