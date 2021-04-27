package controllers

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/controllers/auth"
	"github.com/degary/learn-cmdb/forms"
	"github.com/degary/learn-cmdb/models/k8s"
	k8svc "github.com/degary/learn-cmdb/services/k8s"
	"github.com/degary/learn-cmdb/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sDeployPageController struct {
	LayoutController
}

func (c *K8sDeployPageController) Index() {
	c.Data["menu"] = "k8s_deploy_management"
	c.Data["expand"] = "k8s_management"
	c.TplName = "k8s_deploy_page/index.html"
	c.LayoutSections["LayoutScript"] = "k8s_deploy_page/index.script.html"
}

func (c *K8sDeployPageController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	result, total, queryTotal := k8s.NewDeploymentManager().Query(q, start, length)
	c.Data["json"] = map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"result":          result,
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": queryTotal,
	}
	c.ServeJSON()
}

type K8sDeployController struct {
	auth.LoginRequiredController
}

func (c *K8sDeployController) Delete() {
	pk, _ := c.GetInt("pk")

	deploy, err := k8s.NewDeploymentManager().GetById(pk)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":   500,
			"text":   "删除错误",
			"result": err,
		}
		c.ServeJSON()
	}
	if deploy.Namespace == "kube-system" {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"text": "不能删除kube-system命名空间的deploy",
		}
		c.ServeJSON()
	}
	c.Data["json"] = map[string]interface{}{
		"code":   200,
		"text":   "删除成功",
		"result": pk,
	}

	c.ServeJSON()
}

func (c *K8sDeployController) Modify() {
	//如果是post请求
	if c.Ctx.Input.IsPost() {
		valid := &validation.Validation{}
		form := &forms.K8s_deploy_form{}
		json := map[string]interface{}{
			"code":   400,
			"text":   "提交数据错误",
			"result": nil,
		}
		//如果解析form失败
		if err := c.ParseForm(form); err != nil {
			valid.SetError("error", err.Error())
			json["result"] = valid.Errors
			//如果解析form成功
		} else {
			//如果验证失败
			if ok, err := valid.Valid(form); err != nil {
				fmt.Println("验证失败")
				valid.SetError("error", err.Error())
				json["result"] = valid.Errors
				//如果验证成功
			} else if ok {
				/*
					验证成功后,开始操作数据 更改deployment参数
				*/
				k8sCli := k8svc.NewClient(beego.AppConfig.String("k8sconfig"))
				deployCli, deploy, err := k8sCli.GetDeployment(form.Name, form.Namespace)
				if err != nil {
					json = map[string]interface{}{
						"code":   400,
						"text":   "deploy失败",
						"result": err.Error(),
					}
				} else {
					deploy.Spec.Replicas = utils.Int32ToPtr(form.Replicas)
					deploy.Namespace = form.Namespace
					deploy.Name = form.Name
					deploy.Spec.Template.Spec.Containers[0].Image = form.Image
					deployCli.Update(context.TODO(), deploy, metav1.UpdateOptions{})
					json = map[string]interface{}{
						"code":   200,
						"text":   "创建成功",
						"result": form,
					}
				}
			} else if !ok {
				json = map[string]interface{}{
					"code":   400,
					"text":   "参数验证失败",
					"result": valid.Errors,
				}
			}
		}
		c.Data["json"] = json
		c.ServeJSON()
		//如果是get请求
	} else {
		pk, _ := c.GetInt("pk")
		deploy, _ := k8s.NewDeploymentManager().GetById(pk)
		c.Data["object"] = deploy
		c.TplName = "k8s_deployment/modify.html"
	}

}
