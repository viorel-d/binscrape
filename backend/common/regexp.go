package common

import (
	"regexp"
)

func RegexpFindAll(re string, target []byte) [][]byte {
	rexp := regexp.MustCompile(re)
	result := rexp.FindAll(target, -1)
	return result
}

func RegexpFindAllString(re string, target string) []string {
	rexp := regexp.MustCompile(re)
	var result []string
	partialResult := rexp.FindAllStringSubmatch(target, -1)
	for _, match := range partialResult {
		result = append(result, match[1])
	}
	return result
}
