package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct {
	HTTPClient *http.Client
}

func (s *Service) Download(url, path string) error {
	// Create file
	dir := filepath.Dir(path)
	if _, statErr := os.Stat(dir); statErr != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Do request
	rsp, err := s.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := rsp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err := io.Copy(file, rsp.Body); err != nil {
		return err
	}

	return nil
}
