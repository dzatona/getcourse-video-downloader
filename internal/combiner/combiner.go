package combiner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func CombineFiles(outputFile string) error {
	listFile, err := os.Create("ffmpeg_list_*.txt")
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Error removing file: %v", err)
		}
	}(listFile.Name())

	tsFiles, _ := filepath.Glob("parts/*.ts")
	binFiles, _ := filepath.Glob("parts/*.bin")
	files := append(tsFiles, binFiles...)

	sort.Strings(files)

	for _, file := range files {
		_, err := listFile.WriteString(fmt.Sprintf("file '%s'\n", file))
		if err != nil {
			return err
		}
	}
	err = listFile.Close()
	if err != nil {
		log.Printf("Error closing file: %v", err)
	}

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile.Name(),
		"-c", "copy",
		outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, string(output))
	}

	return nil
}
