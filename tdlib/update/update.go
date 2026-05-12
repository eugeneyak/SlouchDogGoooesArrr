package update

import (
	"encoding/json"
	"errors"
)

type Update interface{}
type Action interface{}

func Unmarshal(data []byte) (Update, error) {
	var resolver Typed

	if err := json.Unmarshal(data, &resolver); err != nil {
		return nil, err
	}

	switch resolver.Type {
	case "updateAuthorizationState":
		var update UpdateAuthorizationState

		if err := json.Unmarshal(data, &update); err != nil {
			return nil, err
		}

		return update, nil

	default:
		return nil, errors.New("unrecognized update type: " + string(resolver.Type))
	}
}
