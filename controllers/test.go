package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
	"time"
)

type TestPageController struct {
	LayoutController
}

func (c *TestPageController) Index() {
	c.TplName = "test_page/index.html"
}

type TestController struct {
	auth.LoginRequiredController
}

func (c *TestController) Test() {
	c.Data["json"] = time.Now()
	c.ServeJSON()
}
