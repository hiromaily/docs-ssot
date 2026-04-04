package include

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var includePattern = regexp.MustCompile(`<!--\s*@include:\s*(.*?)\s*-->`)

func ProcessFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var result string
	scanner := bufio.NewScanner(file)

	var resultSb22 strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		matches := includePattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			includePath := matches[1]

			content, err := os.ReadFile(includePath)
			if err != nil {
				return "", fmt.Errorf("include error (%s): %w", includePath, err)
			}

			resultSb22.WriteString(string(content) + "\n")
		} else {
			resultSb22.WriteString(line + "\n")
		}
	}
	result += resultSb22.String()

	return result, scanner.Err()
}
