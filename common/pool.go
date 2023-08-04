package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type PoolData struct {
	PoolName   string           `json:"poolName"`
	Pool       map[string]Kudos `json:"pool"`
	PoolStatus string           `json:"poolStatus"`
}

func parsePoolData(data string) PoolData {
	var poolData PoolData
	err := json.Unmarshal([]byte(data), &poolData)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return poolData
}

// FIXME: pool throw error sometimes
func PoolGet(poolId string) PoolData {
	cmd := exec.Command("npx", "@loremlabs/setler@latest", "pool", "get", "--poolId", poolId)
	output, err := cmd.Output()
	if err != nil {
		// outputStr := string(output)
		log.Panicln(output)
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// Skip the first line
	if len(lines) > 0 {
		lines = lines[1:]
	}

	modifiedOutput := strings.Join(lines, "\n")

	return parsePoolData(modifiedOutput)
}
