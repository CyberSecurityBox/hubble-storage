package slack

type BlockPayload struct {
	Blocks []blocks `json:"blocks"`
}

type blocks struct {
	Type   string   `json:"type,omitempty"`
	Text   *text     `json:"text,omitempty"`
	Fields []fields `json:"fields,omitempty"`
}

type text struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type fields struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
