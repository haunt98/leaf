package main

import (
	"leaf/internal/pkg/image"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
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

	r := mux.NewRouter()
	r.HandleFunc("/images", imageHandler.HandleReceive).Methods("POST")
	r.HandleFunc("/images/{uuid}", imageHandler.HandleGetStatus).Methods("GET")

	log.Fatal(http.ListenAndServe(":42069", r))
}
