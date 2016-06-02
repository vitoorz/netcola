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
	if strings.Contains(in, "mongow") {
		return taskMongoCreate
	}
	if strings.Contains(in, "env") {
		return taskEnv
	}
	if strings.Contains(in, "service") {
		return taskServiceList
	}
	return taskUnknown
}
