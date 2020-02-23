package copier

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Copy(fromPath, toPath string) error {
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
