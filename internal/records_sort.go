package internal

import (
	"sort"
	"strconv"
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
	tsi, erri := strconv.Atoi(s[i]["ts"])
	tsj, errj := strconv.Atoi(s[j]["ts"])
	if erri != nil && errj == nil {
		return false
	}
	if erri == nil && errj != nil {
		return true
	}
	return tsi < tsj
}

// SortRecords sort records
func SortRecords(result []map[string]string) []map[string]string {
	sort.Sort(records(result))
	return result
}
