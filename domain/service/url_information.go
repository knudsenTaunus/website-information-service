package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/knudsenTaunus/website-information-service/domain/entity"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	wg                                                    sync.WaitGroup
	externalLinkMutex, internalLinkMutex, headingMapMutex sync.RWMutex
)

// WebsiteInformation to retrieve information about a website
type WebsiteInformation struct {
}

// New is a constructor function to provide a new WebsiteInformation Service
func New() *WebsiteInformation {
	return &WebsiteInformation{}
}

// GetWebsiteInformation provides information about the specified website
func (u *WebsiteInformation) GetWebsiteInformation(document io.ReadCloser) (*entity.WebsiteInformation, error) {
	if document == nil {
		return nil, errors.New("failed to create information service - missing website")
	}
	information := u.tokanizeHTML(document)
	brokenLinksChan := make(chan bool)

	for link := range information.ExternalLinks {
		wg.Add(1)
		go u.checkLinkAccessability(link, brokenLinksChan)
	}

	for range information.ExternalLinks {
		if <-brokenLinksChan {
			information.BrokenLinks++
		}
	}
	wg.Done()
	wg.Wait()
	close(brokenLinksChan)
	return information, nil
}

func (u *WebsiteInformation) tokanizeHTML(document io.ReadCloser) *entity.WebsiteInformation {
	tokenizer := html.NewTokenizer(document)
	result := entity.NewWebsiteInformation()
	getTitle := false

	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		switch tt {
		case html.ErrorToken:
			return result
		case html.DoctypeToken:
			if token.Data == "html" {
				result.DoctypeVersion = "HTML 5"
			}
			if strings.Contains(token.Data, "HTML 4.01") {
				result.DoctypeVersion = "HTML 4"
			}
		case html.StartTagToken:
			switch token.DataAtom {
			case atom.A:
				for _, element := range token.Attr {
					if element.Key == "href" {
						wg.Add(1)
						go u.registerLink(element.Val, result)
					}
				}
			case atom.Head:
				getTitle = true

			case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
				result.Headings[token.DataAtom.String()]++
			}

		case html.TextToken:
			if getTitle {
				data := strings.TrimSpace(token.Data)
				if len(data) > 0 {
					result.Title = data
					getTitle = false
				}
			}
		}
	}

}

func (u *WebsiteInformation) registerLink(link string, result *entity.WebsiteInformation) {
	defer wg.Done()
	if strings.HasPrefix(link, "http") {
		externalLinkMutex.Lock()
		result.ExternalLinks[link]++
		externalLinkMutex.Unlock()
		return
	}
	if strings.HasPrefix(link, "#") {
		internalLinkMutex.Lock()
		result.InternalLinks[link]++
		internalLinkMutex.Unlock()
		return
	}

}

func (u *WebsiteInformation) checkLinkAccessability(link string, c chan bool) {
	defer wg.Done()
	_, err := http.Get(link)
	if err != nil {
		c <- true
	}
	c <- false
	return
}
