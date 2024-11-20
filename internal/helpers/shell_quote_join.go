package helpers

import (
	"strconv"
	"strings"
)

// ShellQuoteJoin quotes each string for the shell and joins using the separator
func ShellQuoteJoin(ss []string) string {
	var qs []string
	for _, s := range ss {
		qs = append(qs, strconv.Quote(s))
	}
	return strings.Join(qs, ", ")
}
