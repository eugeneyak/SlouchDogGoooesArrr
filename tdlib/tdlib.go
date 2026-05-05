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
	"fmt"
	"unsafe"
)

type Update struct {
	Type    string          `json:"@type"`
	Payload json.RawMessage `json:"values"`
}

// Client wraps a TDLib JSON client instance.
type Client struct {
	ptr     unsafe.Pointer
	Updates chan Update
}

// Creates a new TDLib JSON client instance.
func Init() *Client {
	return &Client{
		ptr:     C.td_json_client_create(),
		Updates: make(chan Update),
	}
}

// Send sends a JSON-formatted request to the TDLib client.
func (c *Client) Send(request string) {
	if c == nil || c.ptr == nil {
		return
	}
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))
	C.td_json_client_send(c.ptr, cRequest)
}

// Execute synchronously executes a TDLib request.
// The returned string is valid until the next call to Receive or Execute.
func (c *Client) Execute(request string) string {
	if c == nil {
		return ""
	}
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
	if c == nil || c.ptr == nil {
		panic("F")
	}

	for {
		raw := C.td_json_client_receive(c.ptr, C.double(10))
		if raw == nil {
			continue
		}

		str := C.GoString(raw)

		fmt.Println(str)

		var update Update

		if err := json.Unmarshal([]byte(str), &update); err == nil {
			c.Updates <- update
		} else {
			panic("Failed to unmarshal update: " + err.Error())
		}
	}
}

// Destroy releases the TDLib client instance.
func (c *Client) Destroy() {
	if c == nil || c.ptr == nil {
		return
	}
	C.td_json_client_destroy(c.ptr)
	c.ptr = nil
}
