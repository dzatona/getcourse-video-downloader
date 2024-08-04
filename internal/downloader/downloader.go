package downloader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadPlaylist(playlistURL string) ([]string, error) {
	resp, err := http.Get(playlistURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playlist []string
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "http") {
			playlist = append(playlist, line)
		} else if !strings.HasPrefix(line, "#") {
			baseURL, _ := url.Parse(playlistURL)
			relativeURL, _ := url.Parse(line)
			absoluteURL := baseURL.ResolveReference(relativeURL)
			playlist = append(playlist, absoluteURL.String())
		}
	}

	if len(playlist) == 0 {
		scanner = bufio.NewScanner(strings.NewReader(string(body)))
		var lastLine string
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if strings.HasPrefix(lastLine, "http") {
			return DownloadPlaylist(lastLine)
		}
		return nil, fmt.Errorf("no valid URLs found in the playlist")
	}

	return playlist, scanner.Err()
}

func DownloadFiles(playlist []string) error {
	bar := progressbar.NewOptions(len(playlist)-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for i := 1; i < len(playlist); i++ {
		urlString := playlist[i]
		fileName := filepath.Join("parts", fmt.Sprintf("file_%05d.ts", i-1))

		err := downloadFile(urlString, fileName)
		if err != nil {
			return fmt.Errorf("error downloading file %d: %w", i, err)
		}

		err = bar.Add(1)
		if err != nil {
			log.Printf("Error adding bar item to download progress: %v", err)
		}
	}

	fmt.Println()
	return nil
}

func downloadFile(urlString, fileName string) error {
	resp, err := http.Get(urlString)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
