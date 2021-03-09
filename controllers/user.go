package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
	"github.com/degary/learn-cmdb/models"
	"strings"
)

type UserPageController struct {
	LayoutController
}

func (c *UserPageController) Index() {
	c.Data["menu"] = "user_management"
	c.Data["expand"] = "system_management"
	c.LayoutSections["LayoutScript"] = "user_page/index.script.html"
	c.TplName = "user_page/index.html"
}

type UserController struct {
	auth.LoginRequiredController
}

func (c *UserController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	//[]*User,total,queryTotal
	users, total, queryTotal := models.DefaultUserManager.Query(q, start, length)
	c.Data["json"] = map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"result":          users,
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": queryTotal,
	}
	c.ServeJSON()
}

func (c *UserController) Create() {
	if c.Ctx.Input.IsPost() {

	} else {
		c.TplName = "user/create.html"
	}
}
