package fileingestion

import "time"

type IngestionFileName struct {
	period   time.Time
	category string
}

func newIngestionFileName(p time.Time, c string) *IngestionFileName {
	return &IngestionFileName{
		period:   p,
		category: c,
	}
}

func (f *IngestionFileName) Period() time.Time {
	return f.period
}

func (f *IngestionFileName) Category() string {
	return f.category
}
