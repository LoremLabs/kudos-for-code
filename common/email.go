package common

import (
	"log"
	"sync"

	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier()

type ValidationResult struct {
	Email   string
	IsValid bool
}

func ValidateEmails(emails []string) []ValidationResult {
	emailChannel := make(chan string, len(emails))
	resultChannel := make(chan ValidationResult, len(emails))
	var wg sync.WaitGroup
	numWorkers := 20

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for email := range emailChannel {
				isValid := checkEmail(email)
				resultChannel <- ValidationResult{Email: email, IsValid: isValid}
			}
		}()
	}

	for _, email := range emails {
		emailChannel <- email
	}

	close(emailChannel)
	wg.Wait()
	close(resultChannel)

	results := make([]ValidationResult, 0, len(emails))
	for result := range resultChannel {
		results = append(results, result)
	}

	return results
}

func checkEmail(email string) bool {
	ret, err := verifier.Verify(email)
	if err != nil {
		log.Printf("Failed to verify email address %s: %v", email, err)
		return false
	}
	if !ret.Syntax.Valid {
		log.Printf("Invalid email address syntax: %s", email)
		return false
	}

	return true
}
