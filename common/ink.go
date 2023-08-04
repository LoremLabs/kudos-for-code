package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func Ink(poolId string, filePath string, poolEndpoint string) {
	cmd := exec.Command("npx", "@loremlabs/setler@latest", "pool", "ink", "--poolId", poolId, "--inFile", filePath, "--poolEndpoint", poolEndpoint)
	output, err := cmd.Output()
	if err != nil {
		log.Panicln("EEEEE", err)
		panic(err)
	}

	if !strings.Contains(string(output), "✅ Pool inked") {
		log.Panicf("Error: during processing %s", filePath)
	}
}

func ProcessNDJSON(reader *bufio.Reader, chunkSize int) ([]string, []Kudos) {
	lineCount := 0
	var tempFilePaths []string // To store temporary file paths

	var lines []string
	var kudosSlice []Kudos

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading input:", err)
			panic(err)
		}

		lineCount++

		// just test valid json
		var kudos Kudos
		err = json.Unmarshal([]byte(line), &kudos)
		if err != nil {
			fmt.Println("Error parsing NDJSON:", err)
			panic(err)
		}

		kudosSlice = append(kudosSlice, kudos)
		lines = append(lines, line)

		if lineCount%chunkSize == 0 {
			tempFilePath, err := createTempFile(lines)
			if err != nil {
				fmt.Println("Error creating temporary file:", err)
				panic(err)
			}
			tempFilePaths = append(tempFilePaths, tempFilePath)
			lines = []string{}
		}
	}

	if len(lines) > 0 {
		tempFilePath, err := createTempFile(lines)
		if err != nil {
			fmt.Println("Error creating temporary file:", err)
			panic(err)
		}
		tempFilePaths = append(tempFilePaths, tempFilePath)
	}

	return tempFilePaths, kudosSlice
}

func createTempFile(lines []string) (string, error) {
	tempDir := os.TempDir()

	tempFile, err := os.CreateTemp(tempDir, "tempdata_*.ndjson")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	content := []byte(strings.Join(lines, ""))
	_, err = tempFile.Write(content)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func Compare(kudosSlice []Kudos, target []Kudos) []Kudos {
	var difference []Kudos
	for _, sourceKudos := range kudosSlice {
		found := false
		for i, targetKudos := range target {
			if sourceKudos.Id == targetKudos.Id && reflect.DeepEqual(sourceKudos, targetKudos) {
				found = true
				target = append(target[:i], target[i+1:]...)
				break
			}
		}
		if !found {
			difference = append(difference, sourceKudos)
		}
	}

	return difference
}
