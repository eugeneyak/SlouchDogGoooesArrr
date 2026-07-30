package tdlib

/*
#cgo CFLAGS: -I/usr/src/td/tdlib/include
#cgo LDFLAGS: -L/usr/src/td/build -Wl,-rpath,/usr/src/td/build -ltdjson -ldl

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

const (
	Fatal   = 0
	Error   = 1
	Warning = 2
	Info    = 3
	Debug   = 4
	Verbose = 5
)

// TDLib wraps a TDLib JSON client instance.
type TDLib struct {
	ptr unsafe.Pointer
}

// Creates a new TDLib JSON client instance.
func Init() (*TDLib, error) {
	ptr := C.td_json_client_create()

	log.Println("TDLib client initialized:", ptr)

	td := TDLib{
		ptr: unsafe.Pointer(ptr),
	}

	td.SetLogVerbosityLevel(Fatal)

	return &td, nil
}

// Send sends a JSON-formatted request to the TDLib client.
// func (td *TDLib) Send(method string, payload Payload) error {
// 	payload["@type"] = method

// 	json, err := json.Marshal(payload)
// 	if err != nil {
// 		td.log.Println("Error marshaling payload:", err)
// 		return err
// 	}

// 	cRequest := C.CString(string(json))
// 	defer C.free(unsafe.Pointer(cRequest))

// 	C.td_json_client_send(td.ptr, cRequest)
// 	return nil
// }

// Execute synchronously executes a TDLib request.
// The returned string is valid until the next call to Receive or Execute.
func (td *TDLib) execute(method string, payload map[string]any) error {
	payload["@type"] = method

	json, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	CJson := C.CString(string(json))
	defer C.free(unsafe.Pointer(CJson))

	C.td_json_client_execute(td.ptr, CJson)

	return nil
}

// Receive waits for a TDLib response or update for up to timeout seconds.
// It returns the JSON response string or an empty string on timeout.
func (td *TDLib) ReceiveAsync(ctx context.Context) chan []byte {
	channel := make(chan []byte)
	go td.receive(ctx, channel)

	return channel
}

func (td *TDLib) receive(ctx context.Context, updates chan []byte) {
	for {
		select {
		case <-ctx.Done():
			close(updates)
			return

		default:
			raw := C.td_json_client_receive(td.ptr, C.double(2))
			if raw == nil {
				continue
			}

			updates <- []byte(C.GoString(raw))
		}
	}
}

// Destroy releases the TDLib client instance.
func (td *TDLib) Destroy() {
	C.td_json_client_destroy(td.ptr)
	log.Println("TDLib client destroed:", td.ptr)
	td.ptr = nil
}

func (td *TDLib) SetLogVerbosityLevel(level uint8) {
	td.execute("setLogVerbosityLevel", map[string]any{
		"new_verbosity_level": level,
	})
}
