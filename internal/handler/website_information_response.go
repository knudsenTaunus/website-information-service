package handler

// WebsiteInformationResponse Response entity to provide website information
type WebsiteInformationResponse struct {
	DoctypeVersion string         `json:"DoctypeVersion"`
	Title          string         `json:"Title"`
	Headings       map[string]int `json:"Headings"`
	ExternalLinks  int            `json:"ExternalLinks"`
	InternalLinks  int            `json:"InternalLinks"`
	BrokenLinks    int            `json:"BrokenLinks"`
}
