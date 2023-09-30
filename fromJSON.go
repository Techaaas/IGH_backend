// Content struct with Items as []string
type Content struct {
	File  string   `json:"file"`
	Items []string `json:"items"`
}

// ...

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

		var fileContent Content
		fileContent.File = fileDiff.Name

		i := 0
		for _, line := range strings.Split(fullCode, "\n") {
			if deletion, ok := deletions[i+1]; ok {
				fileContent.Items = append(fileContent.Items, fmt.Sprintf("-%s", deletion))
			}

			if addition, ok := additions[i+1]; ok {
				fileContent.Items = append(fileContent.Items, fmt.Sprintf("+%s", addition))
			}

			if _, ok := deletions[i+1]; !ok && _, ok := additions[i+1]; !ok {
				fileContent.Items = append(fileContent.Items, fmt.Sprintf("*%s", line))
			}

			i++
		}

		// Append the file content directly to the result array
		resultArray = append(resultArray, fileContent.File)
		resultArray = append(resultArray, fileContent.Items...)

		fmt.Printf("Output for file %s processed\n", fileDiff.Name)
	}

	// Print or use the result array
	db.addDiffData(resultArray);
}
