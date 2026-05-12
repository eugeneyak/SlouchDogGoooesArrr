package slouchdog

import (
	"fmt"

	"slouchdog/tdlib"
	"slouchdog/tdlib/action"
	"slouchdog/tdlib/update"
)

func Authorize(td *tdlib.TDLib, update update.UpdateAuthorizationState) {
	switch update.AuthorizationState.Type {
	case "authorizationStateWaitTdlibParameters":
		td.Send(action.SetTdlibParameters{
			Type:               "setTdlibParameters",
			APIID:              0,
			APIHash:            "00000000000000000000000000000000",
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
