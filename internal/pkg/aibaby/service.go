package aibaby

import (
	"io"
	"leaf/pkg/downloader"
	"leaf/pkg/imgur"
	"log"
	"os"
	"path/filepath"
)

type Service struct {
	DownloaderService *downloader.Service
	ImgurService      *imgur.Service
}

func (s *Service) DoMagic(url, fromPath, toPath string) (MagicResponse, error) {
	if err := s.DownloaderService.Download(url, fromPath); err != nil {
		return MagicResponse{}, err
	}

	if err := s.copy(fromPath, toPath); err != nil {
		return MagicResponse{}, err
	}

	imgurUploadImgRsp, err := s.ImgurService.UploadImage(toPath)
	if err != nil {
		return MagicResponse{}, err
	}

	if !imgurUploadImgRsp.Success {
		log.Println("Hmm imgur sometimes")
	}

	return MagicResponse{
		URL: imgurUploadImgRsp.Data.Link,
	}, nil
}

func (s *Service) copy(fromPath, toPath string) error {
	// Open file
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := fromFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Create file
	toDir := filepath.Dir(toPath)
	if _, statErr := os.Stat(toDir); statErr != nil {
		if err := os.MkdirAll(toDir, os.ModePerm); err != nil {
			return err
		}
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := toFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err := io.Copy(toFile, fromFile); err != nil {
		return err
	}

	return nil
}
