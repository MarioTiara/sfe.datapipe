package salesfe

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var SalesFERegex = regexp.MustCompile(
	`^SalesSFE_(\d{6})$`,
)

func NewSalesFEPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(SalesFERegex)
}
