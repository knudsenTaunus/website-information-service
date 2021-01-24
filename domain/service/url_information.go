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
	// Remove this; package-scoped WaitGroup is bound to fail because of scoping issues
	//
	//wg                                                    sync.WaitGroup 
	
	// Remove these as well -- these are only called when data within entity.WebsiteInformation is changed. A Mutex *within* entity.WebsiteInformation
	// should guard/control how that struct is locked/unlocked, and then you create methods to update the specific data field (in your code below, a
	// map of some sort) that handles locking then unlocking the field within those methods.
	// This implementation leaves the responsibility of locking/unlocking to the consumer of the entity.WebsiteInformation package/API, which is brittle
	// 
	// externalLinkMutex, internalLinkMutex, headingMapMutex sync.RWMutex   
	
)

// WebsiteInformation to retrieve information about a website

// If this struct has no internal data or resources, why is it a struct? 
// Using pointer receivers for all of these will also create more allocations than needed since, again, no internal data is being used here (currently)
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
	information := u.tokenizeHTML(document)
	brokenLinksChan := make(chan bool)
	
	// If you're using a method in this loop, you'll need a pointer here to ensure wg within this scope
	// is updated when the new goroutine called wg.Done()
	wg: = &sync.WaitGroup{}  

	for link := range information.ExternalLinks {
		wg.Add(1)
		go u.checkLinkAccessability(wg, link, brokenLinksChan)
	}

	for range information.ExternalLinks {
		if <-brokenLinksChan {
			information.BrokenLinks++
		}
	}
	// wg.Done() // Allow created goroutines above control over when they're calling wg.Done(), as above
	wg.Wait()
	close(brokenLinksChan)
	return information, nil
}

func (u *WebsiteInformation) tokenizeHTML(document io.ReadCloser) *entity.WebsiteInformation {
	tokenizer := html.NewTokenizer(document)
	result := entity.NewWebsiteInformation()
	getTitle := false

	wg := &sync.WaitGroup{}
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
						go u.registerLink(wg, element.Val, result)
					}
				}
			case atom.Head:
				getTitle = true

			case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
				// Change this to some sort of a method in the entity.NewWebsiteInformation package
				// That method will implement the locking/unlocking
				// IE this is replaced with (or something similar):
				//  result.IncrementHeading(token.DataAtom.String())
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

func (u *WebsiteInformation) registerLink(wg *sync.WaitGroup, link string, result *entity.WebsiteInformation) {
	defer wg.Done()
	if strings.HasPrefix(link, "http") {
		// externalLinkMutex.Lock()
		// replace with something akin to:
		//     result.IncrementExternalLink(link)
		// Again, the result struct itself has an internal Mutex that is used *within* the IncrementExternalLink
		// to lock/unlock the resource.
		result.ExternalLinks[link]++
		// externalLinkMutex.Unlock()
		return
	}
	if strings.HasPrefix(link, "#") {
		// internalLinkMutex.Lock()
		// replace:
		// result.IncrementInternalLinks(link)
		result.InternalLinks[link]++
		// internalLinkMutex.Unlock()
		return
	}

}

func (u *WebsiteInformation) checkLinkAccessability(wg *sync.WaitGroup, link string, c chan bool) {
	defer wg.Done()
	_, err := http.Get(link)
	if err != nil {
		c <- true
	}
	c <- false
	return
}
