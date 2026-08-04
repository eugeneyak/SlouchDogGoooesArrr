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
	"iter"
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
func (td *TDLib) SendJSON(json string) error {
	cRequest := C.CString(string(json))
	defer C.free(unsafe.Pointer(cRequest))

	C.td_json_client_send(td.ptr, cRequest)
	return nil
}

// Execute synchronously executes a TDLib request.
// The returned string is valid until the next call to Receive or Execute.
func (td *TDLib) execute(method string, payload map[string]any) error {
	payload["@type"] = method

	json, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	return td.executeJSON(string(json))
}

func (td *TDLib) executeJSON(json string) error {
	cJson := C.CString(string(json))
	defer C.free(unsafe.Pointer(cJson))

	C.td_json_client_execute(td.ptr, cJson)
	return nil
}

func (td *TDLib) Updates(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				raw := C.td_json_client_receive(td.ptr, C.double(2))

				if raw == nil {
					continue
				}

				s := C.GoString(raw)

				if !yield(s) {
					return
				}
			}
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
