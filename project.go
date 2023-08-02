package main

import (
	"fmt"
	"log"
	"net/mail"
	"strings"
)

type Project struct {
	name                    string
	dependencies            []*Dependency
	limitDepth              int
	maxDepth                int
	numValidUniqueEmails    int
	numValidEmails          int
	numEmails               int
	numCommits              int
	numCommitsByValidEmails int
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

func NewProject(projectName string, a *AnalyzerResult, limitDepth int) *Project {
	weightFactors := []float32{1, 0.5, 0.25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	p := &Project{
		name:                    projectName,
		dependencies:            []*Dependency{},
		limitDepth:              limitDepth,
		maxDepth:                1,
		numValidUniqueEmails:    0,
		numValidEmails:          0,
		numEmails:               0,
		numCommits:              0,
		numCommitsByValidEmails: 0,
	}

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
		for currentDepth := 1; currentDepth <= p.limitDepth; currentDepth++ {
			p.maxDepth = Max(currentDepth, p.maxDepth)
			nextNodes := map[string]bool{}
			for currentNode := range currentNodes {
				for _, edge := range dg.Edges {
					from := nodeLookup[edge.From]
					to := nodeLookup[edge.To]
					if currentNode == from {
						nextNodes[to] = true
						if depth, ok := packageDepthLookup[from]; ok {
							packageDepthLookup[from] = Min(depth, currentDepth)
						} else {
							packageDepthLookup[from] = currentDepth
						}

						if depth, ok := packageDepthLookup[to]; ok {
							packageDepthLookup[to] = Min(depth, currentDepth)
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
			depth:        limitDepth,
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

		p.numValidUniqueEmails = len(validEmailLookup)
	}

	for _, d := range p.dependencies {
		totalCommits := 0
		for _, c := range d.contributors {
			p.numEmails += 1
			p.numCommits += c.numCommits
			if c.isValidEmail {
				p.numValidEmails += 1
				p.numCommitsByValidEmails += c.numCommits
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
}

func (p *Project) LogProjectStat() {
	log.Printf("== BEGIN:Project Stat ==================\n")
	log.Printf("num dependencies: %d\n", len(p.dependencies))
	log.Printf("limit depth: %d\n", p.limitDepth)
	log.Printf("max depth: %d\n", p.maxDepth)
	log.Printf("num valid unique emails: %d\n", p.numValidUniqueEmails)
	log.Printf("num valid emails: %d\n", p.numValidEmails)
	log.Printf("num emails: %d\n", p.numEmails)
	log.Printf("num commits: %d\n", p.numCommits)
	log.Printf("num commits by valid emails: %d\n", p.numCommitsByValidEmails)
	log.Printf("== END:Project Stat   ==================\n")
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
