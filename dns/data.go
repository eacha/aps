package dns

type OpenResolver struct {
	Questions          []Question `json:"questions,omitempty"`
	Answers            []Answer   `json:"answer,omitempty"`
	RecursionAvailable bool       `json:"recursion_available"`
	ResolveCorrectly   bool       `json:"resolve_correctly"`
	RawResponse        []byte     `json:"raw_response,omitempty"`
}

type DNSData struct {
	IP           string        `json:"ip,omitempty"`
	Error        string        `json:"error,omitempty"`
	OpenResolver *OpenResolver `json:"open_resolveer,omitempty"`
}
