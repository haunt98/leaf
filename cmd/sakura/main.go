package main

import (
	"leaf/pkg/downloader"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("go run download url path")
		return
	}

	httpClient := &http.Client{}

	downloaderService := &downloader.Service{
		HTTPClient: httpClient,
	}

	switch os.Args[1] {
	case "download":
		if err := downloaderService.Download(os.Args[2], os.Args[3]); err != nil {
			log.Println(err)
		}
	}
}
