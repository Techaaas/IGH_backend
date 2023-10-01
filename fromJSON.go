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
	Commit1 string     `json:"commit1"`
	Commit2 string     `json:"commit2"`
	Files   []FileDiff `json:"files"`
}

type Content struct {
	File  string `json:"file"`
	Items []struct {
		Type   string `json:"type"`
		String string `json:"string"`
	} `json:"items"`
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

func main7() {
	db.connector()
	db.dropTables()
	filePath := "diff.json"

	diffData, err := readJSONFromFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Slice to store elements of the result
	var resultArray []string

	// Add commit1 as string
	resultArray = append(resultArray, diffData.Commit1)
	// Add commit2 as string
	resultArray = append(resultArray, diffData.Commit2)

	// Create a map for the Content structure
	contentMap := make(map[string]interface{})

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

		var fileContent Content
		fileContent.File = fileDiff.Name

		i := 0
		for _, line := range strings.Split(fullCode, "\n") {
			if deletion, ok := deletions[i+1]; ok {
				fileContent.Items = append(fileContent.Items, struct {
					Type   string `json:"type"`
					String string `json:"string"`
				}{Type: "-", String: deletion})

				if addition, ok := additions[i+1]; ok {
					fileContent.Items = append(fileContent.Items, struct {
						Type   string `json:"type"`
						String string `json:"string"`
					}{Type: "+", String: addition})

					i++
					continue
				}
			} else if addition, ok := additions[i+1]; ok {
				fileContent.Items = append(fileContent.Items, struct {
					Type   string `json:"type"`
					String string `json:"string"`
				}{Type: "+", String: addition})

				i++
				continue
			} else {
				fileContent.Items = append(fileContent.Items, struct {
					Type   string `json:"type"`
					String string `json:"string"`
				}{Type: "*", String: line})
			}
			i++
		}

		// Marshal the fileContent to a JSON string
		fileContentJSON, err := json.Marshal(fileContent)
		if err != nil {
			fmt.Println("Error marshaling file content to JSON:", err)
			return
		}

		// Add the JSON string to the map
		contentMap[fileDiff.Name] = string(fileContentJSON)

		fmt.Printf("Output for file %s processed\n", fileDiff.Name)
	}

	// Convert the contentMap to a JSON string
	contentJSON, err := json.Marshal(contentMap)
	if err != nil {
		fmt.Println("Error marshaling content to JSON:", err)
		return
	}

	// Add the contentJSON as a string to the resultArray
	resultArray = append(resultArray, string(contentJSON))

	fmt.Println(resultArray[0])
	fmt.Println(resultArray[1])

	// Print the result instead of saving to a file
	db.addDiffData(resultArray)
}
