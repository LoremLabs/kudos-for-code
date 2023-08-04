package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/LoremLabs/kudos-for-code-action/common"
	"github.com/spf13/cobra"
)

var poolId string
var poolEndpoint string
var inkCmd = &cobra.Command{
	Use:     "ink",
	Aliases: []string{"gen"},
	Short:   "Ink kudos from pipe",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		filePaths, kudosSlice := common.ProcessNDJSON(reader, 1000)
		for _, filePath := range filePaths {
			common.Ink(poolId, filePath, poolEndpoint)
		}

		for i, filePath := range filePaths {
			log.Printf("inking: %d/%d\n", i+1, len(filePaths))
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("Error removing temporary file:", err)
			}
		}

		log.Printf("createdFiles: %d", len(filePaths))

		poolData := common.PoolGet(poolId)
		var target []common.Kudos
		for _, kudos := range poolData.Pool {
			target = append(target, kudos)
		}

		diff := common.Compare(kudosSlice, target)

		if len(diff) > 0 {
			fmt.Print("There are some differences between the input and the pool.\n")
			for d := range diff {
				fmt.Println(d)
			}
		} else {
			fmt.Print("The input and the pool are identical.\n")
			fmt.Printf("%d inked\n", len(kudosSlice))
		}
	},
}

func init() {
	inkCmd.Flags().StringVarP(&poolId, "poolId", "i", "", "Pool ID")
	inkCmd.Flags().StringVarP(&poolEndpoint, "poolEndpoint", "e", "", "Pool Endpoint")
	inkCmd.MarkFlagRequired("poolId")

	rootCmd.AddCommand(inkCmd)
}
