package common

import (
	"log"
	"net"
	"net/mail"
	"strings"
	"sync"
)

type ValidationResult struct {
	Email   string
	IsValid bool
}

func ValidateEmails(emails []string, numWorkers int) []ValidationResult {
	emailChannel := make(chan string, len(emails))
	resultChannel := make(chan ValidationResult, len(emails))
	var wg sync.WaitGroup

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
	// Check syntax
	_, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("Invalid email address syntax: %s", email)
		return false
	}

	// Check MX record
	domain := strings.Split(email, "@")[1]
	mx, err := net.LookupMX(domain)
	if err != nil || len(mx) == 0 {
		log.Printf("Failed to verify email address %s: Mail server does not exist: %v", email, err)
		return false
	}

	return true
}
