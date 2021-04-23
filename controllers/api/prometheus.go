package api

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/degary/learn-cmdb/controllers/base"
	"github.com/degary/learn-cmdb/forms"
	"github.com/degary/learn-cmdb/models"
)

type PrometheusController struct {
	base.ApiController
}

func (c *PrometheusController) Alert() {
	c.Ctx.Input.CopyBody(1024 * 1024)
	form := forms.AlertsForm{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &form); err == nil {
		for _, alert := range form.Alerts {
			models.DefaultAlertManager.Notify(alert.ToModel())
		}
	} else {
		logs.Error(err)
	}
	data := map[string]string{
		"code": "200",
	}
	c.Data["json"] = data
	c.ServeJSON()
}
