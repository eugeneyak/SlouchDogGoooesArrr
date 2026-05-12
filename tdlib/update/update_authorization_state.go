package update

type AuthorizationState struct {
	Type string `json:"@type"`
	Link string `json:"link"`
}

type UpdateAuthorizationState struct {
	AuthorizationState AuthorizationState `json:"authorization_state"`
}
