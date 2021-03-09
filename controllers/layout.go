package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
)

type LayoutController struct {
	auth.LoginRequiredController
}

func (c *LayoutController) Prepare() {
	c.LoginRequiredController.Prepare()

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["LayoutStyle"] = ""
	c.LayoutSections["LayoutScript"] = ""

	c.Layout = "layouts/base.html"
	c.Data["menu"] = ""
	c.Data["expand"] = ""
}
