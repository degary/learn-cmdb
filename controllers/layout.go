package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
)

type LayoutController struct {
	auth.LoginRequiredController
}

func (c *LayoutController) Prepare() {
	c.LoginRequiredController.Prepare()
	c.Layout = "layouts/base.html"
}
