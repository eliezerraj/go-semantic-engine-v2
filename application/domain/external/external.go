package external

type VectorSearchRequest struct {
	Text string `json:"text" validate:"required"`
}

type VectorSearchResponse struct {
	Response    string	`json:"response"`
	Vector		any	`json:"vector,omitempty"`
}

type IntentDecomposeRequest struct {
	Intent string `json:"intent" validate:"required"`
}

type IntentDecomposeResponse struct {
	Response    string	`json:"response"`
	Intent		any	`json:"intent,omitempty"`
}
