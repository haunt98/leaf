package image

const (
	Successful = 1
	Processing = 2
)

const (
	Topic = "imageTopic"
	Group = "imageGroup"
)

type Status struct {
	Status       int    `json:"status"`
	OriginalURL  string `json:"original_url"`
	ProcessedURL string `json:"processed_url"`
}

type ProcessMessage struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}

type ReceiveRequest struct {
	URL string `json:"url"`
}

type ReceiveResponse struct {
	UUID string `json:"uuid"`
}
