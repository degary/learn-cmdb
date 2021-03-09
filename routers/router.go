package routers

import (
	"github.com/astaxie/beego"
	"github.com/degary/learn-cmdb/controllers"
	"github.com/degary/learn-cmdb/controllers/auth"
)

func init() {
	beego.AutoRouter(&auth.AuthController{})
	beego.AutoRouter(&controllers.UserPageController{})
	beego.AutoRouter(&controllers.UserController{})

}
