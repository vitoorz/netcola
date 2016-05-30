package task

import "strings"

func Parse(in string) int {
	if strings.Contains(in, "exit") {
		return taskExit
	}
	if strings.Contains(in, "now") {
		return taskGetNow
	}
	if strings.Contains(in, "info") {
		return taskInfo
	}
	if strings.Contains(in, "later") {
		return taskLater
	}
	if strings.Contains(in, "help") {
		return taskHelp
	}
	return taskUnknown
}
