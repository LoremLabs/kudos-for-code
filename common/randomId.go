package common

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/google/uuid"
)

// Generate URL-safe base64 encoded UUID
// It's same with id=$(uuidgen | xxd -r -p | base64 | sed 's/+/-/g; s/\//_/g; s/=//g') in bash.
// The bash script you have provided generates a random UUID (Universally Unique Identifier), converts it to binary data, then encodes it in base64 format, and finally replaces certain characters to make it URL-safe.
//
// From ChatGPT3:
// Here's a breakdown of each part of the script:
//
// 1. `uuidgen`: This command generates a random UUID in its canonical form (e.g., `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`).
// 2. `xxd -r -p`: This part converts the canonical UUID format into binary data.
// 3. `base64`: The binary data from the previous step is then encoded in base64 format.
// 4. `sed 's/+/-/g; s/\//_/g; s/=//g'`: Finally, this `sed` command replaces the characters `+` with `-`, `/` with `_`, and removes any equal signs `=` from the base64-encoded string to make it URL-safe.
//
// The variable `id` will hold the resulting URL-safe base64-encoded UUID.
// Keep in mind that UUIDs are generally used to create unique identifiers,
// and the specific transformations in this script appear to be for the purpose
// of creating a URL-safe version of the UUID (by replacing characters
// that have special meanings in URLs).
func NewRandomId() string {
	// Generate a random UUID
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error generating UUID:", err)
		panic(err)
	}

	// Convert the UUID to binary data
	uuidBytes := uuidObj[:]
	// Encode the binary data in base64
	base64Encoded := base64.RawURLEncoding.EncodeToString(uuidBytes)

	// Replace characters to make it URL-safe
	urlSafeBase64 := strings.ReplaceAll(base64Encoded, "+", "-")
	urlSafeBase64 = strings.ReplaceAll(urlSafeBase64, "/", "_")
	urlSafeBase64 = strings.ReplaceAll(urlSafeBase64, "=", "")

	return urlSafeBase64

}
