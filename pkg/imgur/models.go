package imgur

type UploadImageResponse struct {
	Data    UploadImageData `json:"data"`
	Success bool            `json:"success"`
}

type UploadImageData struct {
	Link string `json:"link"`
}
