package aibaby

import (
	"leaf/pkg/copier"
	"leaf/pkg/downloader"
	"leaf/pkg/imgur"
	"log"
)

type Service struct {
	DownloaderService *downloader.Service
	ImgurService      *imgur.Service
}

func (s *Service) DoMagic(url, fromPath, toPath string) (MagicResponse, error) {
	if err := s.DownloaderService.Download(url, fromPath); err != nil {
		return MagicResponse{}, err
	}

	if err := copier.Copy(fromPath, toPath); err != nil {
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
