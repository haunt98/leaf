package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct {
	HTTPClient *http.Client
	ClientID   string
}

const (
	ImageURL   = "https://api.imgur.com/3/image"
	ImageField = "image"
)

func (s *Service) UploadImage(path string) (UploadImageResponse, error) {
	file, err := os.Open(path)
	if err != nil {
		return UploadImageResponse{}, err
	}

	// Build body request
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile(ImageField, filepath.Base(path))
	if _, err := io.Copy(part, file); err != nil {
		return UploadImageResponse{}, err
	}
	if err := writer.Close(); err != nil {
		return UploadImageResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, ImageURL, reqBody)
	if err != nil {
		return UploadImageResponse{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", s.ClientID))

	// Do request
	rsp, err := s.HTTPClient.Do(req)
	if err != nil {
		return UploadImageResponse{}, err
	}
	defer func() {
		if err := rsp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return UploadImageResponse{}, err
	}

	var uploadImgRsp UploadImageResponse
	if err := json.Unmarshal(rspBody, &uploadImgRsp); err != nil {
		return UploadImageResponse{}, err
	}

	return uploadImgRsp, nil
}
