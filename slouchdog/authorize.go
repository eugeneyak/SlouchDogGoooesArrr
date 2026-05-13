package slouchdog

import (
	"fmt"
	"os"
	"strconv"

	"slouchdog/tdlib"
	"slouchdog/tdlib/action"
	"slouchdog/tdlib/update"
)

func Authorize(td *tdlib.TDLib, update update.UpdateAuthorizationState) {
	switch update.AuthorizationState.Type {
	case "authorizationStateWaitTdlibParameters":
		id, err := strconv.Atoi(os.Getenv("APIID"))
		if err != nil {
			panic("Error converting APIID to integer")
		}

		td.Send(action.SetTdlibParameters{
			Type:               "setTdlibParameters",
			APIID:              int32(id),
			APIHash:            os.Getenv("APIHASH"),
			SystemLanguageCode: "en-US",
			DeviceModel:        "Slouchdog",
			ApplicationVersion: "0.0.1",
		})

	case "authorizationStateWaitPhoneNumber":
		td.Send(action.RequestQrCodeAuthentication{
			Type: "requestQrCodeAuthentication",
		})

	case "authorizationStateWaitOtherDeviceConfirmation":
		fmt.Println("Waiting for other device confirmation...")

	default:
		fmt.Println("Unhandled authorization state:", update.AuthorizationState)
	}
}
