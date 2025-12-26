package masteroutlet

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var MasterOutletRegex = regexp.MustCompile(
	`^(\d{6})_Masterlist Outlet$`,
)

func NewMasterOutletPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(MasterOutletRegex)
}
