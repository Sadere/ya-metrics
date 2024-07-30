package common

import (
	"crypto/rand"
	"strings"
)

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

func GenerateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
