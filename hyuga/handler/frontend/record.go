package frontend

import (
	"hyuga/database"
	"hyuga/handler/base"

	"github.com/gin-gonic/gin"
)

func GetRecords(c *gin.Context) {
	var record database.Record
	switch c.Query("type") {
	case "dns":
		record = database.DnsRecord{}
	case "http":
		record = database.HttpRecord{}
	case "jndi":
		record = database.JndiRecord{}
	default:
		base.ReturnError(c, 101)
		return
	}

	list, err := database.GetUserRecordsByUserID(record, base.GetUserID(c), c.Query("filter"))
	if err != nil {
		base.ReturnError(c, 102, err)
		return
	}

	base.ReturnJSON(c, list)
}

func CleanRecords(c *gin.Context) {
	records := []database.Record{
		database.DnsRecord{},
		database.HttpRecord{},
		database.JndiRecord{},
	}
	for _, r := range records {
		if err := database.DeleteRecordsByUserID(r, base.GetUserID(c)); err != nil {
			base.ReturnError(c, 102, err)
			return
		}
	}
	base.ReturnJSON(c, nil)
}
