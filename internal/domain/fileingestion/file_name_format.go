package fileingestion

import (
	"errors"
	"strings"
	"time"
)

type FileNameFormat interface {
	Parse(raw string)
}

type PrefixPeriodFormat struct{}

func (PrefixPeriodFormat) Parse(raw string) (*IngestionFileName, error) {
	parts := strings.Split(raw, "_")
	if len(parts) < 3 {
		return nil, errors.New("invalid file name format")
	}

	period, err := time.Parse("200601", parts[0])
	if err != nil {
		return nil, err
	}

	category := strings.Join(parts[1:], "_")
	return newIngestionFileName(period, category), nil
}

type SuffixPeriodFormat struct{}

func (SuffixPeriodFormat) Parse(raw string) (*IngestionFileName, error) {
	parts := strings.Split(raw, "_")
	if len(parts) < 3 {
		return nil, errors.New("invalid file name format")
	}

	period, err := time.Parse("200601", parts[len(parts)-1])
	if err != nil {
		return nil, err
	}

	category := strings.Join(parts[:len(parts)-1], "_")
	return newIngestionFileName(period, category), nil
}
