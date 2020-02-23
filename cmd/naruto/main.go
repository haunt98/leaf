package main

import (
	"fmt"
	"leaf/internal/pkg/image"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

func main() {
	// Config
	viper.SetConfigName("naruto")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.Get("redis.port")),
	})

	// Kafka
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.bootstrap.servers"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer kafkaProducer.Close()

	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println(ev.TopicPartition.Error)
				}
			}
		}
	}()

	// Services
	imageStatusRepository := image.StatusRepository{
		RedisClient: redisClient,
	}

	imageProcessProducer := image.ProcessProducer{
		KafkaProducer: kafkaProducer,
	}

	imageService := image.Service{
		StatusRepo:      &imageStatusRepository,
		ProcessProducer: &imageProcessProducer,
	}

	imageHandler := image.Handler{
		Service: &imageService,
	}

	// API
	r := mux.NewRouter()
	r.HandleFunc("/images", imageHandler.HandleReceive).Methods("POST")
	r.HandleFunc("/images/{uuid}", imageHandler.HandleGetStatus).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("api.port")), r))
}
