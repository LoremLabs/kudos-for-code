package main

import (
	"sync"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.NewVerifier()
)

type ValidationResult struct {
	Email   string
	IsValid bool
}

func ValidateEmails(emails []string) []ValidationResult {
	emailChannel := make(chan string)
	var wg sync.WaitGroup
	numWorkers := 20

	results := []ValidationResult{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for email := range emailChannel {
				isValid := checkEmail(email)
				results = append(results, ValidationResult{Email: email, IsValid: isValid})
			}
		}()
	}

	for _, email := range emails {
		emailChannel <- email
	}

	close(emailChannel)
	wg.Wait()

	return results
}

func checkEmail(email string) bool {
	ret, err := verifier.Verify(email)
	if err != nil {
		// fmt.Printf("%s: verify email address failed, error is: %+s\n", email, err)
		return false
	}
	if !ret.Syntax.Valid {
		// fmt.Printf("%s: email address syntax is invalid.\n", email)
		return false
	}

	return true
}
