package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func gitDiffToJSON(commit1, commit2, outputFile string) error {
	diffOutput, err := runGitDiff(commit1, commit2)
	if err != nil {
		return fmt.Errorf("error running git diff: %v", err)
	}

	// Skip if the diff output is empty
	if diffOutput == "" {
		fmt.Printf("No differences found between %s and %s. Skipping...\n", commit1, commit2)
		return nil
	}

	fileDiffs, err := parseDiffOutput(diffOutput)
	if err != nil {
		return fmt.Errorf("error parsing git diff output: %v", err)
	}

	diffData := DiffData{
		Commit1: commit1,
		Commit2: commit2,
		Files:   fileDiffs,
	}

	err = writeJSONToFile(diffData, outputFile)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

func runGitDiff(commit1, commit2 string) (string, error) {
	cmd := exec.Command("git", "diff", commit1, commit2)
	diffOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(diffOutput), nil
}

func parseDiffOutput(diffOutput string) ([]FileDiff, error) {
	lines := strings.Split(diffOutput, "\n")

	var fileDiffs []FileDiff
	var currentFileDiff FileDiff
	var oldLineNumber, newLineNumber int

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@@"):
			old, new := extractLineNumbers(line)
			oldLineNumber, newLineNumber = old, new
		case strings.HasPrefix(line, "+++ "):
			if currentFileDiff.Name != "" {
				fileDiffs = append(fileDiffs, currentFileDiff)
			}
			currentFileDiff = FileDiff{Name: strings.TrimSpace(strings.TrimPrefix(line, "+++ b/"))}
		case strings.HasPrefix(line, "---"):

		case strings.HasPrefix(line, "+"):
			newLineNumber++
			content := line[1:]
			change := Change{Type: "addition", OldLineNumber: 0, NewLineNumber: newLineNumber + 2, Content: content}
			currentFileDiff.Changes = append(currentFileDiff.Changes, change)
		case strings.HasPrefix(line, "-"):
			oldLineNumber++
			content := line[1:]
			change := Change{Type: "deletion", OldLineNumber: oldLineNumber + 2, NewLineNumber: 0, Content: content}
			currentFileDiff.Changes = append(currentFileDiff.Changes, change)
		default:
		}
	}

	if currentFileDiff.Name != "" {
		fileDiffs = append(fileDiffs, currentFileDiff)
	}

	return fileDiffs, nil
}

func writeJSONToFile(data interface{}, outputFile string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func extractLineNumbers(header string) (int, int) {
	var oldStart, oldCount, newStart, newCount int
	_, err := fmt.Sscanf(header, "@@ -%d,%d +%d,%d @@", &oldStart, &oldCount, &newStart, &newCount)
	if err != nil {
		return 0, 0
	}
	return oldStart, newStart
}

func parseGitLog(logOutput string) []string {
	// Split the log output into individual commits
	commitStrings := strings.Split(logOutput, "commit ")
	fmt.Println(commitStrings)

	// Filter out empty strings and trim spaces
	var commits []string
	for _, commitString := range commitStrings {
		trimmedCommit := strings.TrimSpace(commitString)
		if trimmedCommit != "" {
			// Extract only the commit hash (assuming it's the first 7 characters)
			commitHash := trimmedCommit[:7]
			commits = append(commits, commitHash)
		}
	}

	return commits
}

func main2() {
	var outputFile string

	outputFile = "diff.json"

	// Run the git log command without the --all option to get commits from all branches
	cmd := exec.Command("git", "log")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the output
	commits := parseGitLog(string(output))

	// Generate JSON files for all combinations of commits
	for i := 0; i < len(commits); i++ {
		for j := i + 1; j < len(commits); j++ {
			commit1 := commits[i]
			commit2 := commits[j]
			outputFileName := fmt.Sprintf(outputFile)
			err := gitDiffToJSON(commit1, commit2, outputFileName)
			if err != nil {
				return
			}
			main7()
		}
	}
}
