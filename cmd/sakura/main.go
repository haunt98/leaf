package main

import (
	"leaf/pkg/downloader"
	"leaf/pkg/imgur"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("wrong use")
		return
	}

	// Config
	viper.SetConfigName("sakura")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return
	}

	// Services
	httpClient := &http.Client{}

	downloaderService := &downloader.Service{
		HTTPClient: httpClient,
	}

	imgurService := &imgur.Service{
		HTTPClient: httpClient,
		ClientID:   viper.GetString("imgur.clientid"),
	}

	switch os.Args[1] {
	case "download":
		if len(os.Args) != 4 {
			log.Println("go run download url path")
			return
		}

		if err := downloaderService.Download(os.Args[2], os.Args[3]); err != nil {
			log.Println(err)
		}
	case "upload_image":
		if len(os.Args) != 3 {
			log.Println("go run upload_image path")
			return
		}

		imgurUploadImageRsp, err := imgurService.UploadImage(os.Args[2])
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%+v\n", imgurUploadImageRsp)
	}
}
