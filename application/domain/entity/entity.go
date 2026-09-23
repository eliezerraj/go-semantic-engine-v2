package entity

type VectorSearch struct {
	Text string `json:"text,omitempty"`
    Embedding []float64 `json:"embedding,omitempty"`
    Capabilities []Capability `json:"capabilities,omitempty"`
}

type Capability struct {
    Text string `json:"text,omitempty"`
	Name string `json:"name,omitempty"`
    Type string `json:"type,omitempty"`
    Endpoint string `json:"endpoint,omitempty"`
    Uri string `json:"uri,omitempty"`
}

type IntentDecompose struct {
	Intent string `json:"intent,omitempty"`
}

type TeiModel struct {
	Inputs string `json:"inputs,omitempty"`
    Normalize bool `json:"normalize,omitempty"`
    Truncate bool `json:"truncate,omitempty"`
    TruncationDirection string `json:"truncation_direction,omitempty"`
    Embedding []float64 `json:"embedding,omitempty"`
}