package routers

import (
	"github.com/astaxie/beego"
	"github.com/degary/learn-cmdb/controllers"
	"github.com/degary/learn-cmdb/controllers/api"
	"github.com/degary/learn-cmdb/controllers/auth"
	"github.com/degary/learn-cmdb/filters"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	//插入过滤函数,用于注册prometheus函数
	beego.InsertFilter("/*", beego.BeforeExec, filters.BeforeExec)
	//false 意思是允许多个filter
	beego.InsertFilter("/*", beego.AfterExec, filters.AfterExec, false)

	beego.AutoRouter(&auth.AuthController{})
	beego.AutoRouter(&controllers.UserPageController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.TokenController{})
	//云平台
	beego.AutoRouter(&controllers.CloudPlatformPageController{})
	beego.AutoRouter(&controllers.CloudPlatformController{})
	//云主机
	beego.AutoRouter(&controllers.VirtualMachineController{})
	beego.AutoRouter(&controllers.VirtualMachinePageController{})

	//k8s平台
	beego.AutoRouter(&controllers.K8sDeployPageController{})
	beego.AutoRouter(&controllers.K8sServicePageController{})
	beego.AutoRouter(&controllers.K8sDeployController{})
	beego.AutoRouter(&controllers.K8sServiceController{})

	//home
	beego.AutoRouter(&controllers.DashboardPageController{})

	//prometheus
	beego.Handler("/metrics/", promhttp.Handler())

	//alert
	beego.AutoRouter(&api.PrometheusController{})
	beego.AutoRouter(&controllers.AlertPageController{})

}
