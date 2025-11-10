package excel

import "strconv"

func safeGet(row []string, index int) string {
	if len(row) > index {
		return row[index]
	}
	return ""
}

func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}
