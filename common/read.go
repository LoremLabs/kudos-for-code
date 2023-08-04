package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// flow
// a. read kudos input from pipe
// b. feed the input to setler ink
// c. get current kudos from semicolons
// d. compare the current kudos with the new kudos
// show me the code

// create a function for reading from pipe
func ReadFromPipe() {
	log.Println("==> -1. ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
	}

	log.Println("==> -2. ")
	ExecuteSetler("5vwzcFmWzgdgwV8uM8tLKV")
}

// execute "npx", "@loremlabs/setler@latest", "pool", "get", "--poolId", poolId
// and return json from stdout
func ExecuteSetler(poolId string) {
	log.Println("==> 1. ")
	cmd := exec.Command("npx", "@loremlabs/setler@latest", "pool", "get", "--poolId", poolId)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the output to string
	outputStr := string(output)

	// Split the output into lines
	lines := strings.Split(outputStr, "\n")

	// Skip the first line
	if len(lines) > 0 {
		lines = lines[1:]
	}

	// Join the lines back into a string
	modifiedOutput := strings.Join(lines, "\n")

	// Assuming the modified output is in []byte format, you can marshal it into JSON
	var result interface{}
	err = json.Unmarshal([]byte(modifiedOutput), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the JSON result
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}
