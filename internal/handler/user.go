package handler

import (
	"hyuga/internal/db"
	"hyuga/internal/handler/base"

	"github.com/gin-gonic/gin"
)

type user struct {
	db *db.Client
}

func NewUser(db *db.Client) *user {
	return &user{
		db: db,
	}
}

func (u *user) Route(e *gin.Engine, middleware ...gin.HandlerFunc) {
	g := e.Group("/user/v2")
	g.POST("/signin", u.signin)
	g.POST("/signup", u.signup)

	auth := g.Group("")
	{
		auth.Use(middleware...)
		auth.POST("/logout", u.logout)
		auth.POST("/change_password", u.changePassword)
	}
}

type signinParams struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (u *user) signin(c *gin.Context) {
	var param = &signinParams{}
	if base.BindValidate(c, param) {
		return
	}

	user, err := u.db.LoginUser(c.Request.Context(),
		param.Username,
		param.Password)
	if err != nil {
		base.ReturnError(c, 2000)
		return
	}

	base.ReturnJSON(c, map[string]interface{}{
		"id":        user.ID,
		"sid":       user.Sid,
		"username":  user.Username,
		"token":     user.Token,
		"code":      user.InviteCode,
		"from_user": user.FromUser,
	})
}

type signupParams struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Code     string `json:"code" form:"code" binding:"required"`
}

func (u *user) signup(c *gin.Context) {
	var param = &signupParams{}
	if base.BindValidate(c, param) {
		return
	}

	_, err := u.db.GreateUser(c.Request.Context(),
		param.Username,
		param.Password,
		param.Code)
	if err != nil {
		base.ReturnError(c, 1002)
		return
	}
	base.ReturnJSON(c, "ok")
}

func (u *user) logout(c *gin.Context) {
	_, err := u.db.LogoutUser(c.Request.Context(), base.GetUserID(c))
	if err != nil {
		base.ReturnError(c, 1002)
		return
	}

	base.ReturnJSON(c, "ok")
}

type changePasswordParams struct {
	Username    string `json:"username" form:"username" binding:"required"`
	OldPassword string `json:"oldpassword" form:"oldpassword" binding:"required"`
	NewPassword string `json:"newpassword" form:"newpassword" binding:"required"`
}

func (u *user) changePassword(c *gin.Context) {
	var param = &changePasswordParams{}
	if base.BindValidate(c, param) {
		return
	}
	_, err := u.db.ChangePasswordUser(c.Request.Context(), param.Username, param.OldPassword, param.NewPassword)
	if err != nil {
		base.ReturnError(c, 2000)
		return
	}

	base.ReturnJSON(c, "ok")
}
