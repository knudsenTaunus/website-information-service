# website-information-service

This is a simple web application to receive some information about a website provided as a URL parameter like so:

 - http://localhost:8080/getWebsiteInformation?website=http://www.spiegel.de

How to run the application:
 - The project uses go mod, so please first 'go get' the packages listed in go.mod
 - Please run the main.go file inside cmd/server - the server will start on port 8080
 - The only endpoint provided is /getWebsiteInformation and takes a URL parameter called website
 - Call the endpoint and provide the website URL that you want to get information from