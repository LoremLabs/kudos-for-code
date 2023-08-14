package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/LoremLabs/kudos-for-code/common"
	"github.com/spf13/cobra"
)

var poolId string
var poolEndpoint string
var validateResult bool
var chunkSize int
var inkCmd = &cobra.Command{
	Use:     "ink",
	Aliases: []string{"gen"},
	Short:   "Ink kudos from pipe",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		filePaths, kudosSlice := common.ProcessNDJSON(reader, chunkSize)

		log.Println("Inking started.")
		for i, filePath := range filePaths {
			var progress int
			if i == len(filePaths)-1 {
				progress = len(kudosSlice)
			} else {
				progress = (i + 1) * chunkSize
			}

			log.Printf("inking: %d/%d\n", progress, len(kudosSlice))
			common.Ink(poolId, filePath, poolEndpoint)
		}

		for _, filePath := range filePaths {
			err := os.Remove(filePath)
			if err != nil {
				log.Println("Error removing temporary file:", err)
			}
		}

		fmt.Printf("%d kudos inked.\n", len(kudosSlice))

		if validateResult {
			log.Println("Validation started.")

			poolData, err := common.PoolGet(poolId)
			if err != nil {
				log.Printf("Validation skipped due to get pool error: %v\n", err)
				fmt.Printf("But, kudos are inked successfully.\n")

				return
			}

			var target []common.Kudos
			for _, kudos := range poolData.Pool {
				target = append(target, kudos)
			}

			diff := common.Compare(kudosSlice, target)

			if len(diff) > 0 {
				fmt.Print("There are some differences between the input and the pool.\n")
				for d := range diff {
					log.Println(d)
				}
			} else {
				fmt.Print("The input and the pool are identical.\n")
			}
		}
	},
}

func init() {
	inkCmd.Flags().StringVarP(&poolId, "poolId", "i", "", "Pool ID")
	inkCmd.Flags().StringVarP(&poolEndpoint, "poolEndpoint", "e", "https://api.semicolons.com", "Pool Endpoint")
	inkCmd.Flags().BoolVarP(&validateResult, "validateResult", "v", false, "Validate the result")
	inkCmd.Flags().IntVarP(&chunkSize, "chunkSize", "c", 1000, "Chunk size for ink request")
	inkCmd.MarkFlagRequired("poolId")

	rootCmd.AddCommand(inkCmd)
}
