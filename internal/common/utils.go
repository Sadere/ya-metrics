package common

import "strings"

func BuildInfo(buildVersion string, buildDate string, buildCommit string) string {
	var sb strings.Builder

	writePart := func(msg string, val string) {
		sb.WriteString(msg)

		if len(val) > 0 {
			sb.WriteString(val)
		} else {
			sb.WriteString("N/A")
		}
		sb.WriteString("\n")
	}

	writePart("Build version: ", buildVersion)
	writePart("Build date: ", buildDate)
	writePart("Build commit: ", buildCommit)

	return sb.String()
}
