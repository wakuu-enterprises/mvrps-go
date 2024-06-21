package video

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ProcessSegments(segmentsDir, outputDir string) {
	// Ensure the output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	// Create a temporary file for the ffmpeg concat list
	concatFilePath := filepath.Join(outputDir, "concat_list.txt")
	files, err := ioutil.ReadDir(segmentsDir)
	if err != nil {
		fmt.Printf("Error reading segments directory: %v\n", err)
		return
	}

	var concatList strings.Builder
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp4" {
			concatList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(segmentsDir, file.Name())))
		}
	}

	err = ioutil.WriteFile(concatFilePath, []byte(concatList.String()), os.ModePerm)
	if err != nil {
		fmt.Printf("Error writing concat list file: %v\n", err)
		return
	}

	// Command to concatenate and structure the segments
	outputFilePath := filepath.Join(outputDir, "output.mp4")
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", concatFilePath, "-c", "copy", outputFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error processing segments: %v\n", err)
		return
	}

	fmt.Printf("Segments processed: %s\n", outputFilePath)
	// Clean up temporary concat list file
	os.Remove(concatFilePath)
}
