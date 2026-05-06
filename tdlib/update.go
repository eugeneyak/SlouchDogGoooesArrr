package tdlib

import (
	"encoding/json"
	"errors"
	"slouchdog/tdlib/update"
)

type Update interface {
	Handle()
}

func Unmarshal(data []byte) (Update, error) {
	var resolver update.Typed

	if err := json.Unmarshal(data, &resolver); err != nil {
		return nil, err
	}

	switch resolver.Type {
	case "updateAuthorizationState":
		var update update.UpdateAuthorizationState

		if err := json.Unmarshal(data, &update); err != nil {
			return nil, err
		}

		return update, nil

	default:
		return nil, errors.New("unrecognized update type: " + string(resolver.Type))
	}
}
