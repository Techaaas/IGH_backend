package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Change struct {
	Type          string `json:"type"`
	OldLineNumber int    `json:"old_line_number"`
	NewLineNumber int    `json:"new_line_number"`
	Content       string `json:"content"`
}

type FileDiff struct {
	Name    string   `json:"name"`
	Changes []Change `json:"changes"`
}

type DiffData struct {
	Commit1 string      `json:"commit1"`
	Commit2 string      `json:"commit2"`
	Files   []FileDiff `json:"files"`
}

func gitDiffToJSON(commit1, commit2, outputFile string) error {
	diffOutput, err := runGitDiff(commit1, commit2)
	if err != nil {
		return fmt.Errorf("error running git diff: %v", err)
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

func main() {
	var commit1, commit2, outputFile string

	fmt.Print("Enter the hash of the first commit: ")
	fmt.Scan(&commit1)

	fmt.Print("Enter the hash of the second commit: ")
	fmt.Scan(&commit2)

	fmt.Print("Enter the name of the output JSON file: ")
	fmt.Scan(&outputFile)

	err := gitDiffToJSON(commit1, commit2, outputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Differences saved to %s\n", outputFile)
}
