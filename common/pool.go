package common

import (
	"encoding/json"
	"fmt"
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

func PoolGet(poolId string) (PoolData, error) {
	cmd := exec.Command("npx", "@loremlabs/setler@latest", "pool", "get", "--poolId", poolId)
	output, err := cmd.Output()
	if err != nil {
		return PoolData{}, fmt.Errorf("error running command: %w: %s", err, output)
	}

	return parsePoolData(output), nil
}
