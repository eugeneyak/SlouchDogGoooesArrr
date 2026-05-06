package update

type Type string

type Typed struct {
	Type Type `json:"@type"`
}
