package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func GenerateEmails(repoUrls []string, noMerges bool) map[string][]string {
	var lookup sync.Map

	//   1: 8.34min
	//   3: 2.59min
	//  10: 59sec
	//  15: 40sec
	//  20: 38sec
	//  25: 34sec
	//  30: 49sec
	// 100: 55sec
	maxConcurrent := 20
	semaphore := make(chan struct{}, maxConcurrent)

	for _, repoUrl := range repoUrls {
		semaphore <- struct{}{} // Acquire a semaphore slot

		go func(url string) {
			defer func() { <-semaphore }() // Release the semaphore slot when done

			destPath, err := os.MkdirTemp("", "dir")
			if err != nil {
				log.Fatal(err)
			}

			emails := []string{}
			if err := cloneRepository(url, destPath); err != nil {
				// Just skip if error
				// e.g., https://github.com/substack/node-concat-map.git
				// log.Printf("Error cloning repository %s: %s\n", url, err)
			} else {
				emails = getAuthorEmails(destPath, noMerges)
			}

			defer os.RemoveAll(destPath)

			// An empty emails means error on git clone(most case)
			// TODO: do better for error case
			lookup.Store(url, emails)
		}(repoUrl)
	}

	// Wait for all goroutines to finish
	for i := 0; i < maxConcurrent; i++ {
		semaphore <- struct{}{}
	}

	normalMap := make(map[string][]string)
	lookup.Range(func(key, value interface{}) bool {
		normalMap[key.(string)] = value.([]string)
		return true
	})

	return normalMap
}

func cloneRepository(repoURL, destPath string) error {
	cmd := exec.Command("git", "clone", "--filter=tree:0", repoURL, destPath)
	err := cmd.Run()
	if err != nil {
		// https://github.com/substack/text-table is not exist so can't be cloned
		return fmt.Errorf("error cloning repository %s: %w", repoURL, err)
	}

	return nil
}

func getAuthorEmails(destPath string, noMerges bool) []string {
	args := []string{"log", "--format='%ae'"}
	if noMerges {
		args = append(args, "--no-merges")
	}
	// --no-merges
	cmd := exec.Command("git", args...)
	cmd.Dir = destPath
	out, _ := cmd.Output()

	cleanedEmails := []string{}
	for _, e := range strings.Split(string(out), "\n") {
		cleanedEmails = append(cleanedEmails, strings.ReplaceAll(e, "'", ""))
	}

	return cleanedEmails
}
