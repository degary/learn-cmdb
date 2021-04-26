package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
	"github.com/degary/learn-cmdb/models/k8s"
	"strings"
)

type K8sDeployPageController struct {
	LayoutController
}

func (c *K8sDeployPageController) Index() {
	c.Data["menu"] = "k8s_deploy_management"
	c.Data["expand"] = "k8s_management"
	c.TplName = "k8s_deploy_page/index.html"
	c.LayoutSections["LayoutScript"] = "k8s_deploy_page/index.script.html"
}

func (c *K8sDeployPageController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	result, total, queryTotal := k8s.NewDeploymentManager().Query(q, start, length)
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

type K8sDeployController struct {
	auth.LoginRequiredController
}

func (c *K8sDeployController) Delete() {
	pk, _ := c.GetInt("pk")

	deploy, err := k8s.NewDeploymentManager().GetById(pk)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":   500,
			"text":   "删除错误",
			"result": err,
		}
		c.ServeJSON()
	}
	if deploy.Namespace == "kube-system" {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"text": "不能删除kube-system命名空间的deploy",
		}
		c.ServeJSON()
	}
	c.Data["json"] = map[string]interface{}{
		"code":   200,
		"text":   "删除成功",
		"result": pk,
	}

	c.ServeJSON()
}

func (c *K8sDeployController) Modify() {
	if c.Ctx.Input.IsPost() {

	} else {
		pk, _ := c.GetInt("pk")
		deploy, _ := k8s.NewDeploymentManager().GetById(pk)
		c.Data["object"] = deploy
		c.TplName = "k8s_deployment/modify.html"
	}

}
