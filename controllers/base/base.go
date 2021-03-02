package base

import (
	"fmt"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	fmt.Println("xsrf")
	c.Data["xsrf_token"] = c.XSRFToken()
}
