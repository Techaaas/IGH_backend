package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
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

type Content struct {
	File  string   `json:"file"`
	Items []string `json:"items"`
}

func readJSONFromFile(filePath string) (DiffData, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return DiffData{}, err
	}

	var diffData DiffData
	err = json.Unmarshal(fileContent, &diffData)
	if err != nil {
		return DiffData{}, err
	}

	return diffData, nil
}

func getFullCode(commit, fileName string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", commit, fileName))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func main() {
    filePath := "diff.json"

    diffData, err := readJSONFromFile(filePath)
    if err != nil {
        fmt.Println("Error reading JSON file:", err)
        return
    }

    // Array to store the required elements
    var resultArray []string

    // Add commit1 and commit2
    resultArray = append(resultArray, diffData.Commit1, diffData.Commit2)

    for _, fileDiff := range diffData.Files {
        additions := make(map[int]string)
        deletions := make(map[int]string)

        for _, change := range fileDiff.Changes {
            if change.Type == "addition" {
                additions[change.NewLineNumber] = change.Content
            } else if change.Type == "deletion" {
                deletions[change.OldLineNumber] = change.Content
            }
        }

        var additionKeys []int
        var deletionKeys []int

        for key := range additions {
            additionKeys = append(additionKeys, key)
        }
        sort.Ints(additionKeys)

        for key := range deletions {
            deletionKeys = append(deletionKeys, key)
        }
        sort.Ints(deletionKeys)

        fullCode, err := getFullCode(diffData.Commit2, fileDiff.Name)
        if err != nil {
            fmt.Printf("Error getting full code for file %s: %v\n", fileDiff.Name, err)
            continue
        }

        resultArray = append(resultArray, fileDiff.Name)

        i := 0
        for _, line := range strings.Split(fullCode, "\n") {
            if deletion, ok := deletions[i+1]; ok {
                resultArray = append(resultArray, fmt.Sprintf("-%s", deletion))
            }

            if addition, ok := additions[i+1]; ok {
                resultArray = append(resultArray, fmt.Sprintf("+%s", addition))
            }

            if _, ok := deletions[i+1]; !ok {
                if _, ok := additions[i+1]; !ok {
                    resultArray = append(resultArray, fmt.Sprintf("*%s", line))
                }
            }

            i++
        }

        fmt.Printf("Output for file %s processed\n", fileDiff.Name)
    }

    // Print or use the result array
    db.addDiffData(resultArray)
}
