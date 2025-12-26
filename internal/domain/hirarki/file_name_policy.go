package hirarki

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var HirarkiRegex = regexp.MustCompile(
	`^(\d{6})_Hirarki$`,
)

func NewHirarkiPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(HirarkiRegex)
}
