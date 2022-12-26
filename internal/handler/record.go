package handler

import (
	"hyuga/internal/db"
	"hyuga/internal/handler/base"

	"github.com/gin-gonic/gin"
)

type record struct {
	db *db.Client
}

func NewRecord(db *db.Client) *record {
	return &record{
		db: db,
	}
}

func (r *record) Route(e *gin.Engine, middleware ...gin.HandlerFunc) {
	g := e.Group("/record/v2")
	g.Use(middleware...)
	g.GET("/search", r.search)
}

type searchParams struct {
	Type    string `json:"type" form:"type"`
	Keyword string `json:"keyword" form:"keyword"`
}

func (r *record) search(c *gin.Context) {
	var param = &searchParams{}
	if base.BindValidate(c, param) {
		return
	}

	records, err := r.db.SearchRecord(c.Request.Context(), base.GetUserID(c), param.Type, param.Keyword)
	if err != nil {
		base.ReturnError(c, 1002)
		return
	}

	base.ReturnJSON(c, map[string]any{
		"total": len(records),
		"list":  records,
	})
}
