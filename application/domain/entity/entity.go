package entity

type VectorSearch struct {
	Text string `json:"text,omitempty"`
    Capabilities []Capability `json:"capabilities,omitempty"`
}

type Capability struct {
	Name string `json:"name,omitempty"`
    Type string `json:"type,omitempty"`
    Endpoint string `json:"endpoint,omitempty"`
    Uri string `json:"uri,omitempty"`
}

type IntentDecompose struct {
	Intent string `json:"intent,omitempty"`
}