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
	"slouchdog/tdlib/update"
	"unsafe"
)

// TDLib wraps a TDLib JSON client instance.
type TDLib struct {
	ptr unsafe.Pointer
	log *log.Logger
}

// Creates a new TDLib JSON client instance.
func Init() *TDLib {
	log := log.New(log.Writer(), "[TDLib] ", log.LstdFlags|log.Lshortfile)
	ptr := C.td_json_client_create()

	log.Println("TDLib client initialized:", ptr)

	return &TDLib{
		ptr: unsafe.Pointer(ptr),
		log: log,
	}
}

// Send sends a JSON-formatted request to the TDLib client.
func (td *TDLib) Send(action update.Action) error {
	json, err := json.Marshal(action)
	if err != nil {
		td.log.Println("Error marshaling action:", err)
		return err
	}

	cRequest := C.CString(string(json))
	defer C.free(unsafe.Pointer(cRequest))

	C.td_json_client_send(td.ptr, cRequest)
	return nil
}

// Execute synchronously executes a TDLib request.
// The returned string is valid until the next call to Receive or Execute.
func (td *TDLib) Execute(request string) string {
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))

	result := C.td_json_client_execute(td.ptr, cRequest)
	if result == nil {
		return ""
	}

	return C.GoString(result)
}

// Receive waits for a TDLib response or update for up to timeout seconds.
// It returns the JSON response string or an empty string on timeout.
func (td *TDLib) Receive(ctx context.Context) chan update.Update {
	channel := make(chan update.Update)
	go td.receive(ctx, channel)

	return channel
}

func (td *TDLib) receive(ctx context.Context, updates chan update.Update) {
	for {
		select {
		case <-ctx.Done():
			td.log.Println("Receive loop exiting due to context cancellation")
			td.log.Println(context.Cause(ctx))
			close(updates)
			return

		default:
			raw := C.td_json_client_receive(td.ptr, C.double(10))
			if raw == nil {
				continue
			}

			json := C.GoString(raw)

			td.log.Println(json)

			update, error := update.Unmarshal([]byte(json))

			if error == nil {
				updates <- update
			}
		}
	}
}

// Destroy releases the TDLib client instance.
func (td *TDLib) Destroy() {
	C.td_json_client_destroy(td.ptr)
	td.ptr = nil
}
