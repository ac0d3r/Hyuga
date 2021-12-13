package frontend

import (
	"hyuga/database"

	"github.com/gin-gonic/gin"
)

func GetRecords(c *gin.Context) {
	var record database.Record
	switch c.Query("type") {
	case "dns":
		record = database.DnsRecord{}
	case "http":
		record = database.HttpRecord{}
	default:
		ReturnError(c, 101)
		return
	}

	list, err := database.GetUserRecordsByUserID(record, UserID(c), c.Query("filter"))
	if err != nil {
		ReturnError(c, 102, err)
		return
	}

	ReturnJSON(c, list)
}

func CleanRecords(c *gin.Context) {
	records := []database.Record{
		database.DnsRecord{},
		database.HttpRecord{},
	}
	for _, r := range records {
		if err := database.DeleteRecordsByUserID(r, UserID(c)); err != nil {
			ReturnError(c, 102, err)
			return
		}
	}
	ReturnJSON(c, nil)
}
