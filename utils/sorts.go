package utils

import (
	"sort"
	"strconv"

	"github.com/labstack/gommon/log"
)

type records []map[string]string

func (s records) Len() int {
	return len(s)
}

func (s records) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s records) Less(i, j int) bool {
	// 处理 ts
	tsi, err := strconv.Atoi(s[i]["ts"])
	if err != nil {
		log.Debug(err)
		return true
	}
	tsj, err := strconv.Atoi(s[j]["ts"])
	if err != nil {
		log.Debug(err)
		return true
	}
	return tsi < tsj
}

// SortRecords sort records
func SortRecords(result []map[string]string) []map[string]string {
	sort.Sort(records(result))
	return result
}
