package controllers

import (
	"github.com/degary/learn-cmdb/models/k8s"
	"strings"
)

type K8sServicePageController struct {
	LayoutController
}

func (c *K8sServicePageController) Index() {
	c.Data["menu"] = "k8s_service_management"
	c.Data["expand"] = "k8s_management"
	c.TplName = "k8s_service_page/index.html"
	c.LayoutSections["LayoutScript"] = "k8s_service_page/index.script.html"
}

func (c *K8sServicePageController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	result, total, queryTotal := k8s.NewServicePortManager().Query(q, start, length)
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

func (c *K8sServicePageController) Create() {

}
