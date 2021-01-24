package main

import (
	"github.com/knudsenTaunus/website-information-service/internal/handler"
	"github.com/knudsenTaunus/website-information-service/internal/server"
)

func main() {

	websiteInformationHandler := handler.New()
	server := server.New(websiteInformationHandler)
	server.Start()
}
