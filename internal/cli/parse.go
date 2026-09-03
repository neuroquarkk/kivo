package cli

import "strings"

func parseArgs(line string) []string {
	return strings.Fields(line)
}
