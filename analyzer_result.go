package main

import (
	"encoding/json"
	"log"
	"os"
)

// Package represents a package in the JSON data
type Package struct {
	ID               string   `json:"id"`
	PURL             string   `json:"purl"`
	Authors          []string `json:"authors"`
	DeclaredLicenses []string `json:"declared_licenses"`
	Description      string   `json:"description"`
	HomepageURL      string   `json:"homepage_url"`
	BinaryArtifact   struct {
		URL  string `json:"url"`
		Hash struct {
			Value     string `json:"value"`
			Algorithm string `json:"algorithm"`
		} `json:"hash"`
	} `json:"binary_artifact"`
	SourceArtifact struct {
		URL  string `json:"url"`
		Hash struct {
			Value     string `json:"value"`
			Algorithm string `json:"algorithm"`
		} `json:"hash"`
	} `json:"source_artifact"`
	VCS struct {
		Type     string `json:"type"`
		URL      string `json:"url"`
		Revision string `json:"revision"`
		Path     string `json:"path"`
	} `json:"vcs"`
	VCSProcessed struct {
		Type     string `json:"type"`
		URL      string `json:"url"`
		Revision string `json:"revision"`
		Path     string `json:"path"`
	} `json:"vcs_processed"`
}

// DependencyGraph represents a dependency graph in the JSON data
type DependencyGraph struct {
	Packages []string `json:"packages"`
	Scopes   map[string][]struct {
		Root int `json:"root"`
	} `json:"scopes"`
	Nodes []struct {
		PackageIndex int `json:"pkg"`
	} `json:"nodes"`
	Edges []struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"edges"`
}

// AnalyzerResult represents the root of the JSON data
type AnalyzerResult struct {
	Analyzer struct {
		Result struct {
			Packages         []Package                  `json:"packages"`
			DependencyGraphs map[string]DependencyGraph `json:"dependency_graphs"`
		} `json:"result"`
	} `json:"analyzer"`
}

func NewAnalyzerResult(filePath string) *AnalyzerResult {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Create an instance of AnalyzerResult to hold the parsed JSON data
	var analyzerResult AnalyzerResult

	// Unmarshal the JSON data into the AnalyzerResult struct
	err = json.Unmarshal(jsonData, &analyzerResult)
	if err != nil {
		log.Println("Error unmarshaling JSON data:", err)
		panic(err)
	}

	return &analyzerResult
}
