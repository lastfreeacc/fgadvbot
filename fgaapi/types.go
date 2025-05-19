package fgaapi

// AdviceResponse represents response wrapper
type AdviceResponse struct {
	Status string   `json:"status"`
	Errors []string `json:"errors"`
	Data   []Data   `json:"data"`
}

// Data represents advice
type Data struct {
	ID          int64      `json:"id"`
	Text        string     `json:"text"`
	HTML        string     `json:"html"`
	Tags        []string   `json:"tags"`
	Conclusions []struct{} `json:"conclusions"`
}

// Advice uses in api v1
type Advice struct {
	ID    int64  `json:"id"`
	Text  string `json:"text"`
	Sound string `json:"sound"`
}
