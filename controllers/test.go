package controllers

import "github.com/degary/learn-cmdb/controllers/auth"

type TestController struct {
	auth.LoginRequiredController
}

func (c *TestController) Test() {
}
