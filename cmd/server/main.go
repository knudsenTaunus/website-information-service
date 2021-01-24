package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/knudsenTaunus/website-information-service/domain/service"
)

func main() {
	startGetCall := time.Now()
	resp, err := http.Get("https://www.spiegel.de/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	finishGetCall := time.Since(startGetCall)

	startAnalyze := time.Now()

	informationService := service.New()

	websiteInfo, err := informationService.GetWebsiteInformation(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	finishAnalyze := time.Since(startAnalyze)
	fmt.Printf("Title: %s\n", websiteInfo.Title)
	fmt.Printf("Doctype Version: %s\n", websiteInfo.DoctypeVersion)
	fmt.Printf("Internal Links: %d\n", len(websiteInfo.InternalLinks))
	fmt.Printf("External Links: %d\n", len(websiteInfo.ExternalLinks))
	fmt.Printf("Broken Links: %d\n", websiteInfo.BrokenLinks)

	fmt.Println("Headers:")
	for level, amount := range websiteInfo.Headings {
		fmt.Printf("%s: %d\n", level, amount)
	}
	fmt.Printf("HTTP Get took: %s\n", finishGetCall)
	fmt.Printf("Analyze took: %s\n", finishAnalyze)
}
