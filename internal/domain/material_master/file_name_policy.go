package materialmaster

import (
	"regexp"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
)

var MaterialmasterRegex = regexp.MustCompile(
	`^(\d{6})_SFE_Tablemaster_Material Outlet$`,
)

func NewMaterialmasterPolicy() *fileingestion.RegexFileNamePolicy {
	return fileingestion.NewRegexFileNamePolicy(MaterialmasterRegex)
}
