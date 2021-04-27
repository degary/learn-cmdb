package forms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/services/k8s"
	"strings"
)

type K8s_deploy_form struct {
	Name      string `form:"name" json:"name"`
	Namespace string `form:"namespace" json:"namespace"`
	Image     string `form:"image" json:"image"`
	Replicas  int32  `form:"replicas" json:"replicas"`
}

func (c *K8s_deploy_form) Valid(v *validation.Validation) {
	c.Namespace = strings.TrimSpace(c.Namespace)
	c.Name = strings.TrimSpace(c.Name)
	c.Image = strings.TrimSpace(c.Image)

	v.Required(c.Name, "name.name").Message("请填写name")
	v.Min(c.Replicas, 1, "name.name").Message("replicas最小数字为1")
	v.Required(c.Image, "name.name").Message("请填写image")

	client := k8s.NewClient(beego.AppConfig.String("k8sconfig"))
	nsList := client.NamespacesNameList()
	nsMap := map[string]int{}
	for n, ns := range nsList {
		nsMap[ns] = n
	}
	if _, ok := nsMap[c.Namespace]; !ok {
		v.SetError("Namespace", "namespace不存在,请创建")
	}

}
