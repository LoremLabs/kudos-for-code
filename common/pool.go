package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type PoolData struct {
	PoolName   string           `json:"poolName"`
	Pool       map[string]Kudos `json:"pool"`
	PoolStatus string           `json:"poolStatus"`
}

func parsePoolData(data []byte) PoolData {
	var poolData PoolData
	err := json.Unmarshal(data, &poolData)
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

	return parsePoolData(output)
}
