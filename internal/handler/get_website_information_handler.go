package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/knudsenTaunus/website-information-service/internal/domain/entity"
)

// WebsiteInformationService Interface definition for usage inside handler
type WebsiteInformationService func(document io.ReadCloser) (*entity.WebsiteInformation, error)

// New returns a new handlerFunc for website information calls
func New(websiteInformationService WebsiteInformationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		website := r.URL.Query().Get("website")
		resp, err := http.Get(website)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		websiteInformation, err := websiteInformationService(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		websiteInformationResponse := &WebsiteInformationResponse{
			Title:          websiteInformation.Title,
			DoctypeVersion: websiteInformation.DoctypeVersion,
			Headings:       websiteInformation.Headings,
			InternalLinks:  len(websiteInformation.InternalLinks),
			ExternalLinks:  len(websiteInformation.ExternalLinks),
			BrokenLinks:    websiteInformation.BrokenLinks,
		}

		response, err := json.Marshal(websiteInformationResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

}
