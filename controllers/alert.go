package controllers

import (
	"github.com/degary/learn-cmdb/models"
	"strings"
)

type AlertPageController struct {
	LayoutController
}

func (c *AlertPageController) Index() {
	c.Data["menu"] = "alert"
	c.TplName = "alert_page/index.html"
	c.LayoutSections["LayoutScript"] = "alert_page/index.script.html"
}

func (c *AlertPageController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	result, total, queryTotal := models.DefaultAlertManager.Query(q, start, length)
	c.Data["json"] = map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"result":          result,
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": queryTotal,
	}
	c.ServeJSON()
}
