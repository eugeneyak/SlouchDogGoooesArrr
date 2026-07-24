package bridge

type Message struct {
	body    string
	headers []string
}

func newMessage(body string) Message {
	return Message{
		body:    body,
		headers: make([]string, 0),
	}
}

func (m Message) header(name, value string) {
	m.headers = append(m.headers, name, value)
}
