package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/knudsenTaunus/website-information-service/internal/domain/entity"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	externalLinkMutex, internalLinkMutex, headingMapMutex sync.RWMutex
)

// GetWebsiteInformation provides information about the specified website
func GetWebsiteInformation(document io.ReadCloser) (*entity.WebsiteInformation, error) {
	if document == nil {
		return nil, errors.New("failed to create information service - missing website")
	}
	var wg sync.WaitGroup
	information := collectInformation(document)
	brokenLinksChan := make(chan bool, 1)

	for link := range information.ExternalLinks {
		wg.Add(1)
		go checkLinkAccessability(link, brokenLinksChan, &wg)
	}

	for range information.ExternalLinks {
		if <-brokenLinksChan {
			information.BrokenLinks++
		}
	}

	wg.Wait()
	close(brokenLinksChan)
	return information, nil
}

func collectInformation(document io.ReadCloser) *entity.WebsiteInformation {
	var wg sync.WaitGroup
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
						go registerLink(element.Val, result, &wg)
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
		wg.Wait()
	}
}

func registerLink(link string, result *entity.WebsiteInformation, wg *sync.WaitGroup) {
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

func checkLinkAccessability(link string, c chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get(link)
	if err != nil {
		c <- true
		return
	}
	c <- false
	return
}
