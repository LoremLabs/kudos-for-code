package main

import (
	"fmt"
	"os"

	"github.com/LoremLabs/kudos-for-code-action/common"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a project name and a analyzerResultFilePath as an argument.")
		return
	}
	noMerges := true
	onlyValidEmails := true
	limitDepth := 2 // depth should be less than 6
	projectName := os.Args[2]
	analyzerResultFilePath := os.Args[2]

	analyzerResult := common.NewAnalyzerResult(analyzerResultFilePath)
	p := common.NewProject(projectName, analyzerResult, limitDepth)
	p.EnrichContributors(noMerges)
	p.ScoreContributors(onlyValidEmails)
	p.LogProjectStat()

	for _, k := range common.GenerateKudos(p) {
		fmt.Println(string(k.ToJSON()))
	}
}
