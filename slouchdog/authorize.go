package slouchdog

import (
	"fmt"
	"os"
	"strconv"

	"slouchdog/tdlib"
	"slouchdog/tdlib/update"
)

var AuthorizationState string

func Authorize(td *tdlib.TDLib, update update.UpdateAuthorizationState) {
	fmt.Println(update)
	fmt.Println(update.AuthorizationState)

	AuthorizationState = update.AuthorizationState.Type

	switch update.AuthorizationState.Type {
	case "authorizationStateWaitTdlibParameters":
		authorizeStateWaitTdlibParameters(td)

	case "authorizationStateWaitPhoneNumber":
		authorizeStateWaitPhoneNumber(td)

	case "authorizationStateWaitOtherDeviceConfirmation":
		fmt.Println("Waiting for other device confirmation...")

	default:
		fmt.Println("Unhandled authorization state:", update.AuthorizationState)
	}
}

func authorizeStateWaitTdlibParameters(td *tdlib.TDLib) {
	id, err := strconv.Atoi(os.Getenv("APIID"))
	if err != nil {
		panic("Error converting APIID to integer")
	}

	td.SetTdlibParameters(int32(id), os.Getenv("APIHASH"), "en-US", "Slouchdog", "0.0.1")
}

func authorizeStateWaitPhoneNumber(td *tdlib.TDLib) {
	td.RequestQrCodeAuthentication()
}
