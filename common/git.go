package common

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type CommitAuthor struct {
	Name  string
	Email string
}

func GenerateCommitAuthors(repoUrls []string, noMerges bool) map[string][]CommitAuthor {
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

			var commitAuthors []CommitAuthor
			if err := cloneRepository(url, destPath); err != nil {
				// Just skip if error
				// e.g., https://github.com/substack/node-concat-map.git
				// log.Printf("Error cloning repository %s: %s\n", url, err)
			} else {
				commitAuthors = getAuthorEmails(destPath, noMerges)
			}

			defer os.RemoveAll(destPath)

			// An empty emails means error on git clone(most case)
			// TODO: do better for error case
			lookup.Store(url, commitAuthors)
		}(repoUrl)
	}

	// Wait for all goroutines to finish
	for i := 0; i < maxConcurrent; i++ {
		semaphore <- struct{}{}
	}

	normalMap := make(map[string][]CommitAuthor)
	lookup.Range(func(key, value interface{}) bool {
		normalMap[key.(string)] = value.([]CommitAuthor)
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

func getAuthorEmails(destPath string, noMerges bool) []CommitAuthor {
	delimiter := ":+://:+:"
	args := []string{"log", fmt.Sprintf("--format='%%an%s%%ae'", delimiter)}
	if noMerges {
		args = append(args, "--no-merges")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = destPath
	out, _ := cmd.Output()

	var commitAuthors []CommitAuthor
	for _, row := range strings.Split(string(out), "\n") {
		if row == "" {
			continue
		}

		parts := strings.Split(row[1:len(row)-1], delimiter)
		commitAuthor := CommitAuthor{
			"",
			"",
		}

		if len(parts) >= 2 {
			commitAuthor.Name = parts[0]
			commitAuthor.Email = parts[1]
		}

		commitAuthors = append(commitAuthors, commitAuthor)
	}

	return commitAuthors
}
