package controllers

import "github.com/degary/learn-cmdb/models"

type DashboardPageController struct {
	LayoutController
}

func (c *DashboardPageController) Index() {
	//c.LayoutController.Ctx.ResponseWriter.WriteHeader(200)
	c.LayoutController.Ctx.ResponseWriter.Status = 200
	c.Data["menu"] = "dashboard"
	c.Data["virtualMachineCount"] = models.DefaultVirtualMachineManager.Count()
	c.TplName = "dashboard_page/index.html"
	c.LayoutSections["LayoutScript"] = "dashboard_page/index.script.html"

}
