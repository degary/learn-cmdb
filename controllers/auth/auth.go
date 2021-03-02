package auth

import (
	"github.com/degary/learn-cmdb/controllers/base"
	"github.com/degary/learn-cmdb/models"
)

type LoginRequiredController struct {
	base.BaseController
	User *models.User
}

func (c *LoginRequiredController) Prepare() {
	c.BaseController.Prepare()
	//判断是session认证还是token认证
	if user := DefaultManager.IsLogin(c); user == nil {
		//未登录
		DefaultManager.GoToLoginPage(c) //TODO 需要修改参数

		c.StopRun()
	} else {
		//已登录
		c.User = user
		c.Data["user"] = user //传给前端,使前端也可以看到登录信息
	}
}

type AuthController struct {
	base.BaseController
}

func (c *AuthController) Login() {
	DefaultManager.Login(c)

}
func (c *AuthController) Logout() {
	DefaultManager.Logout(c)

}
