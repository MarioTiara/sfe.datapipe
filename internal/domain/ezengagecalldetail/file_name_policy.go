package ezengagecalldetail

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var EzEngagecalldetailRegex = regexp.MustCompile(
	`^(\d{6})_Call Detailed eZEngage$`,
)

func EzEngagecalldetailPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(EzEngagecalldetailRegex)
}
