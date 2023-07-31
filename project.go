package main

import (
	"fmt"
	"log"
	"net/mail"
	"strings"
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

func (p *Project) EnrichContributors(noMerges bool) {
	vcsURLs := []string{}
	for _, d := range p.dependencies {
		if d.vcsType == "Git" {
			vcsURLs = append(vcsURLs, d.vcsUrl)
		}
	}

	vcsUrlEmailsLookup := GenerateEmails(vcsURLs, noMerges)
	for _, d := range p.dependencies {
		numCommitsPerEmail := map[string]int{}
		for _, email := range vcsUrlEmailsLookup[d.vcsUrl] {
			numCommitsPerEmail[email] += 1
		}

		for email, numCommits := range numCommitsPerEmail {
			d.contributors[email] = &Contributor{
				email:        email,
				isValidEmail: false,
				numCommits:   numCommits,
				score:        0,
			}
		}

	}
}

func (p *Project) ScoreContributors(onlyValidEmails bool) {
	if onlyValidEmails {
		emailLookup := map[string]bool{}
		domainLookup := map[string]bool{}
		for _, d := range p.dependencies {
			for _, c := range d.contributors {
				_, err := mail.ParseAddress(c.email)
				if err == nil {
					emailLookup[c.email] = true
				} else {
					log.Printf("SKIP(%s): %s\n", c.email, err)
					continue
				}

				components := strings.Split(c.email, "@")
				domainLookup[components[1]] = false
			}
		}

		testEmails := []string{}
		for domain := range domainLookup {
			testEmails = append(testEmails, fmt.Sprintf("a@%s", domain))
		}

		emailValidationResults := ValidateEmails(testEmails)
		for _, r := range emailValidationResults {
			components := strings.Split(r.Email, "@")
			domainLookup[components[1]] = r.IsValid
		}

		validEmailLookup := map[string]bool{}
		for email := range emailLookup {
			components := strings.Split(email, "@")
			if domainLookup[components[1]] {
				validEmailLookup[email] = true
			}
		}

		for _, d := range p.dependencies {
			for _, c := range d.contributors {
				c.isValidEmail = validEmailLookup[c.email]
			}
		}

		log.Printf("#valid unique emails: %d\n", len(validEmailLookup))
	}

	countEmails := 0
	countValidEmails := 0
	countCommits := 0
	countValidCommits := 0

	for _, d := range p.dependencies {
		totalCommits := 0
		for _, c := range d.contributors {
			// for stat
			countEmails += 1
			countCommits += c.numCommits
			if c.isValidEmail {
				countValidEmails += 1
				countValidCommits += c.numCommits
			}

			// for logic
			if onlyValidEmails && !c.isValidEmail {
				continue
			}

			totalCommits += c.numCommits
		}

		for _, c := range d.contributors {
			if onlyValidEmails && !c.isValidEmail {
				continue
			}

			c.score = float32(c.numCommits) / float32(totalCommits) * d.weight
		}
	}

	log.Printf("countEmails: %d\n", countEmails)
	log.Printf("countValidEmails: %d\n", countValidEmails)
	log.Printf("countCommits: %d\n", countCommits)
	log.Printf("countValidCommits: %d\n", countValidCommits)
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
