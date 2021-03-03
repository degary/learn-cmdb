package controllers

import (
	"github.com/degary/learn-cmdb/controllers/auth"
	"time"
)

type TestController struct {
	auth.LoginRequiredController
}

func (c *TestController) Test() {
	c.Data["json"] = time.Now()
	c.ServeJSON()
}
