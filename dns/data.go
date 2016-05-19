package dns

// OpenResolver is a struct to store the data obtained by dns.Conn.OpenResolver
type OpenResolver struct {
	Questions          []Question `json:"questions,omitempty"`
	Answers            []Answer   `json:"answer,omitempty"`
	RecursionAvailable bool       `json:"recursion_available"`
	ResolveCorrectly   bool       `json:"resolve_correctly"`
	RawResponse        []byte     `json:"raw_response,omitempty"`
}

// Data is a struct to store the data obtained for the differents queries in dns.Scan
type Data struct {
	IP           string        `json:"ip,omitempty"`
	Error        string        `json:"error,omitempty"`
	OpenResolver *OpenResolver `json:"open_resolveer,omitempty"`
}
