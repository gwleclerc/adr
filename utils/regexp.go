package utils

import "regexp"

var recordNamePatternRegex = regexp.MustCompile(`(\d+)_\S+`)

func GetRecordNumber(name string) string {
	match := recordNamePatternRegex.FindStringSubmatch(name)
	if len(match) <= 0 {
		return ""
	}
	return match[1]
}
