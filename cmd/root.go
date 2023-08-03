package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kudos-for-code",
	Short: "Kudos for Code recognizes contributors and their dependencies.",
	Long:  `Kudos for Code generates recognition (kudos) for contributors based on their project's dependencies, using input from ORT's analyzer-result.json. Complete documentation is available at https://github.com/LoremLabs/kudos-for-code`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
