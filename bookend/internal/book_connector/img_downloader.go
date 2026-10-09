package bookconnector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndSaveCover(isbn10 string, isbn13 string, basePath string) (string, error) {
	targetFileName := fmt.Sprintf("%s.jpg", isbn13)
	targetPath := filepath.Join(basePath, targetFileName)

	// 1. Check if cover already exists
	if _, err := os.Stat(targetPath); err == nil {
		fmt.Printf("Cover already exists at %s, skipping download.\n", targetPath)
		return targetPath, nil
	}

	// 2. Try OpenLibrary
	if isbn13 != "" {
		openlibraryThumbnailUrl := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn13)
		path, err := downloadHelper(isbn13, openlibraryThumbnailUrl, basePath)
		if err == nil {
			return path, nil
		}
		fmt.Println("Failed to download from OpenLibrary:", err)
	}

	// 3. Try Amazon Fallback
	if isbn10 != "" {
		amazonThumbnailUrl := fmt.Sprintf("https://images-na.ssl-images-amazon.com/images/P/%s.01.LZZZZZZZ.jpg", isbn10)
		path, err := downloadHelper(isbn13, amazonThumbnailUrl, basePath)
		if err == nil {
			return path, nil
		}
		fmt.Println("Failed to download from Amazon:", err)
	}

	return "", fmt.Errorf("unable to download cover for ISBN13: %s / ISBN10: %s", isbn13, isbn10)
}

func downloadHelper(isbn string, remoteURL string, basePath string) (string, error) {
	// 1. Download image stream
	resp, err := http.Get(remoteURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// 2. Define relative path and local file path
	relPath := fmt.Sprintf("covers/%s.jpg", isbn)
	fullPath := filepath.Join(basePath, relPath)

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		fmt.Printf("Failed to create folder")
		return "", err
	}

	// 3. Save to disk
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Failed to create file")
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return relPath, err
}
