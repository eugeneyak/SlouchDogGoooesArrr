package tdlib

type Type string

type Update struct {
	Type    Type `json:"@type"`
	Payload string
}
