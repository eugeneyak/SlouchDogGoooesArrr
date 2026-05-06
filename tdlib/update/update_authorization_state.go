package update

import (
	"encoding/json"
	"fmt"
)

type WaitTdlibParameters string

type UpdateAuthorizationState struct {
	AuthorizationState WaitTdlibParameters `json:"authorization_state"`
}

func (t *WaitTdlibParameters) UnmarshalJSON(data []byte) error {
	var typed Typed

	if err := json.Unmarshal(data, &typed); err != nil {
		return err
	}

	*t = WaitTdlibParameters(typed.Type)

	return nil
}

func (u UpdateAuthorizationState) Handle() {
	fmt.Println(u.AuthorizationState)
}
