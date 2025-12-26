package customerfe

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var customerSFERegex = regexp.MustCompile(
	`^(\d{6})_SFE_CustomersList$`,
)

func NewSFECustomerListPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(customerSFERegex)
}
