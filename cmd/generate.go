package cmd

import (
	"fmt"
	"log"

	"github.com/LoremLabs/kudos-for-code-action/common"
	"github.com/spf13/cobra"
)

var noMerges bool
var validEmails bool
var limitDepth int
var projectName string
var analyzerResultFilePath string
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate Kudos from ORT result",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := common.NewAnalyzerResult(analyzerResultFilePath)
		if err != nil {
			log.Fatalf("Failed to create analyzer result: %v", err)
		}
		p := common.NewProject(projectName, result, limitDepth)
		p.EnrichContributors(noMerges)
		p.ScoreContributors(validEmails)
		p.LogProjectStat()

		for _, k := range common.GenerateKudos(p) {
			fmt.Println(string(k.ToJSON()))
		}
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&noMerges, "nomerges", "m", false, "Exclude merge commits")
	generateCmd.Flags().BoolVarP(&validEmails, "validemails", "v", true, "Include only valid emails")
	generateCmd.Flags().IntVarP(&limitDepth, "limitdepth", "d", 2, "Limit of dependency depth")
	generateCmd.Flags().StringVarP(&projectName, "projectname", "n", "test-project", "Project name for the result")
	generateCmd.Flags().StringVarP(&analyzerResultFilePath, "inputfilepath", "i", "", "ORT analyzer result file path")

	generateCmd.MarkFlagFilename(analyzerResultFilePath)

	rootCmd.AddCommand(generateCmd)
}
