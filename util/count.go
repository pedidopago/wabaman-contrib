package util

import (
	"regexp"
	"strconv"
)

var (
	variableCountRegexp = regexp.MustCompile(`\{\{([1-9]*[0-9])\}\}`)
)

func CountAndValidateTemplateVariables(textContent string) (count int, valid bool) {
	matches := variableCountRegexp.FindAllStringSubmatch(textContent, -1)
	nextvalidvar := 1
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		intval, _ := strconv.Atoi(match[1])

		if intval != nextvalidvar {
			return count, false
		}

		nextvalidvar++
		count++
	}

	return count, true
}
