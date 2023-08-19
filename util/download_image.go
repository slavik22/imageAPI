package util

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadImage(url, folderPath string) (string, error) {
	_, filename := filepath.Split(url)
	filename = strings.Split(filename, "?")[0]

	path := filepath.Join(folderPath, filename)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, os.ModePerm)

		if err != nil {
			return "", err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Write the image data to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
