package main

import (
	"fmt"
	"leaf/internal/pkg/aibaby"
	"leaf/internal/pkg/image"
	"leaf/pkg/downloader"
	"leaf/pkg/imgur"
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"

	"github.com/spf13/viper"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	log.Println("sasuke running...")

	// Config
	viper.SetConfigName("sasuke")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.Get("redis.port")),
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Println(err)
		return
	}

	// Kafka
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.bootstrap.servers"),
		"group.id":          image.Group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := kafkaConsumer.SubscribeTopics([]string{image.Topic, "^aRegex.*[Tt]opic"}, nil); err != nil {
		log.Println(err)
		return
	}

	// Services
	imageStatusRepository := image.StatusRepository{
		RedisClient: redisClient,
	}

	imageProcessConsumer := image.ProcessConsumer{
		KafkaConsumer: kafkaConsumer,
	}

	httpClient := &http.Client{}

	downloaderService := &downloader.Service{
		HTTPClient: httpClient,
	}

	imgurService := &imgur.Service{
		HTTPClient: httpClient,
		ClientID:   viper.GetString("imgur.clientid"),
	}

	aibabyService := &aibaby.Service{
		DownloaderService: downloaderService,
		ImgurService:      imgurService,
	}

	for {
		processMsg, err := imageProcessConsumer.ReceiveProcess()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Process %+v\n", processMsg)

		status, err := imageStatusRepository.Get(processMsg.UUID)
		if err != nil {
			log.Println(err)
			continue
		}

		if status.OriginalURL != processMsg.URL {
			log.Println("Hmm something wrong, I can feel it.")
			continue
		}

		fromPath := fmt.Sprintf("./input/%s", processMsg.UUID)
		toPath := fmt.Sprintf("./output/%s", processMsg.UUID)
		magicRsp, err := aibabyService.DoMagic(status.OriginalURL, fromPath, toPath)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Magic %+v", magicRsp)

		status.Status = image.Successful
		status.ProcessedURL = magicRsp.URL
		if err := imageStatusRepository.Set(processMsg.UUID, status); err != nil {
			log.Println(err)
			continue
		}

		log.Println("Process LGTM")
	}
}
