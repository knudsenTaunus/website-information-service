package main

import (
	"github.com/knudsenTaunus/website-information-service/internal/domain/service"
	"github.com/knudsenTaunus/website-information-service/internal/handler"
	"github.com/knudsenTaunus/website-information-service/internal/server"
)

func main() {
	websiteInformationService := service.New()
	websiteInformationHandler := handler.New(websiteInformationService)
	server := server.New(websiteInformationHandler)
	server.Start()
}
