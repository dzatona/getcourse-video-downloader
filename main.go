package main

import (
	"log"
	"os"

	"getcourse-video-downloader/internal/combiner"
	"getcourse-video-downloader/internal/downloader"
	"getcourse-video-downloader/internal/utils"
)

func main() {
	utils.ClearScreen()
	utils.PrintWelcomeMessage()
	utils.ClearParts()

	log.SetFlags(log.Ldate | log.Ltime)

	if len(os.Args) < 3 {
		log.Println("Usage: ./getcourse-video-downloader <playlist_url> <output_file>")
		return
	}

	playlistURL := os.Args[1]
	outputFile := os.Args[2]

	err := os.MkdirAll("parts", 0755)
	if err != nil {
		log.Println("Error creating parts directory:", err)
		return
	}

	empty, err := utils.IsDirEmpty("parts")
	if err != nil {
		log.Println("Error checking if parts directory is empty:", err)
		return
	}

	if empty {
		log.Println("Downloading playlist...")
		playlist, err := downloader.DownloadPlaylist(playlistURL)
		if err != nil {
			log.Println("Error downloading playlist:", err)
			return
		}

		log.Printf("Found %d chunks in the playlist.\n", len(playlist)-1)

		err = downloader.DownloadFiles(playlist)
		if err != nil {
			log.Println("Error downloading files:", err)
			return
		}
		log.Println("Files downloaded successfully!")
	} else {
		log.Println("Parts directory is not empty. Skipping download.")
	}

	log.Println("Combining files...")
	err = combiner.CombineFiles(outputFile)
	if err != nil {
		log.Println("Error combining files:", err)
		return
	}
	log.Println("Files combined successfully!")

	utils.ClearParts()
	log.Println("Temporary files cleaned up successfully!")
}
