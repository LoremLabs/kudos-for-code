package main

import (
	"fmt"
)

type Project struct {
	name         string
	dependencies []*Dependency
}

type Dependency struct {
	id           string //package ID
	vcsType      string
	vcsUrl       string
	depth        int
	weight       float32
	contributors map[string]*Contributor
}

type Contributor struct {
	email        string
	isValidEmail bool
	numCommits   int
	score        float32
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func NewProject(projectName string, a *AnalyzerResult) *Project {
	maxDepth := 10
	weightFactors := []float32{1, 0.5, 0.25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	p := new(Project)
	p.name = projectName
	p.dependencies = []*Dependency{}

	packageDepthLookup := map[string]int{}
	for _, dg := range a.Analyzer.Result.DependencyGraphs {
		roots := map[string]bool{}
		for _, s := range dg.Scopes {
			for _, v := range s {
				roots[dg.Packages[v.Root]] = true
			}
		}

		nodeLookup := []string{}
		for _, node := range dg.Nodes {
			nodeLookup = append(nodeLookup, dg.Packages[node.PackageIndex])
		}

		currentNodes := roots
		for currentDepth := 1; currentDepth <= maxDepth; currentDepth++ {
			nextNodes := map[string]bool{}
			for currentNode := range currentNodes {
				for _, edge := range dg.Edges {
					from := nodeLookup[edge.From]
					to := nodeLookup[edge.To]
					if currentNode == from {
						nextNodes[to] = true
						if depth, ok := packageDepthLookup[from]; ok {
							packageDepthLookup[from] = min(depth, currentDepth)
						} else {
							packageDepthLookup[from] = currentDepth
						}

						if depth, ok := packageDepthLookup[to]; ok {
							packageDepthLookup[to] = min(depth, currentDepth)
						} else {
							packageDepthLookup[to] = currentDepth
						}
					}
				}
			}

			if len(nextNodes) == 0 {
				for node := range currentNodes {
					packageDepthLookup[node] = currentDepth
				}

				break
			}

			currentNodes = nextNodes
		}
	}

	for _, pkg := range a.Analyzer.Result.Packages {
		depth, ok := packageDepthLookup[pkg.ID]
		if !ok {
			continue
		}

		d := &Dependency{
			id:           pkg.ID,
			vcsType:      pkg.VCSProcessed.Type,
			vcsUrl:       pkg.VCSProcessed.URL,
			depth:        maxDepth,
			weight:       0,
			contributors: map[string]*Contributor{},
		}

		d.depth = depth
		d.weight = weightFactors[packageDepthLookup[pkg.ID]-1]
		p.dependencies = append(p.dependencies, d)

	}

	return p
}

func (p *Project) EnrichContributors() {
	vcsURLs := []string{}
	for _, d := range p.dependencies {
		if d.vcsType == "Git" {
			vcsURLs = append(vcsURLs, d.vcsUrl)
		}
	}

	vcsUrlEmailsLookup := GenerateEmails(vcsURLs)
	for _, d := range p.dependencies {
		counter := map[string]int{}
		total := 0
		emails := vcsUrlEmailsLookup[d.vcsUrl]
		emailValidationResults := ValidateEmails(emails)

		for _, result := range emailValidationResults {
			if result.IsValid {
				counter[result.Email] += 1
				total += 1
			}

			d.contributors[result.Email] = &Contributor{
				email:        result.Email,
				isValidEmail: result.IsValid,
			}
		}

		for email, numCommits := range counter {
			d.contributors[email].numCommits = numCommits
			d.contributors[email].score = float32(numCommits) / float32(total) * d.weight
		}
	}
}

func (p *Project) ShowDependencyStat() {
	for _, d := range p.dependencies {
		fmt.Println("==")
		fmt.Printf("id: %s\n", d.id)
		fmt.Printf("depth: %d\n", d.depth)
		fmt.Printf("weight: %f\n", d.weight)

		sum := float32(0)
		for _, v := range d.contributors {
			sum += v.score
		}
		fmt.Printf("sum: %f\n", sum)

		fmt.Printf("#contributors: %d\n", len(d.contributors))
	}
}
