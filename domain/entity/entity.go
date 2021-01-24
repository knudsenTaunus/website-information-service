package entity

// WebsiteInformation Entity object for the retrieved website information
type WebsiteInformation struct {
	DoctypeVersion string
	Title          string
	Headings       map[string]int
	ExternalLinks  map[string]int
	InternalLinks  map[string]int
	BrokenLinks    int
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
