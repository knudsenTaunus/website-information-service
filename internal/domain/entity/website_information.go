package entity

// WebsiteInformation Entity object for the retrieved website information
type WebsiteInformation struct {
	DoctypeVersion string         `json:"DoctypeVersion"`
	Title          string         `json:"Title"`
	Headings       map[string]int `json:"Headings"`
	ExternalLinks  map[string]int `json:"ExternalLinks"`
	InternalLinks  map[string]int `json:"InternalLinks"`
	BrokenLinks    int            `json:"BrokenLinks"`
}

// NewWebsiteInformation Constructor for a new entity object for the retrieved website information
func NewWebsiteInformation() *WebsiteInformation {
	return &WebsiteInformation{
		DoctypeVersion: "",
		Title:          "",
		Headings:       make(map[string]int, 0),
		ExternalLinks:  make(map[string]int, 0),
		InternalLinks:  make(map[string]int, 0),
		BrokenLinks:    0,
	}
}
