package forms

import (
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/cloud"
	"github.com/degary/learn-cmdb/models"
	"strings"
)

type CloudPlatformCreateForm struct {
	Name      string `form:"name"`
	Type      string `form:"type"`
	Addr      string `form:"addr"`
	AccessKey string `form:"access_key"`
	SecretKey string `form:"secret_key"`
	Region    string `form:"region"`
	Remark    string `form:"remark"`
}

func (c *CloudPlatformCreateForm) Valid(v *validation.Validation) {
	c.Name = strings.TrimSpace(c.Name)
	c.Type = strings.TrimSpace(c.Type)
	c.Addr = strings.TrimSpace(c.Addr)
	c.AccessKey = strings.TrimSpace(c.AccessKey)
	c.SecretKey = strings.TrimSpace(c.SecretKey)
	c.Remark = strings.TrimSpace(c.Remark)
	c.Region = strings.TrimSpace(c.Region)

	v.AlphaDash(c.Name, "name.name").Message("名字只能由大小写英文,数字,下划线组成")
	v.MinSize(c.Name, 5, "name.name").Message("名字程度必须在%d-%d之内", 5, 32)
	v.MaxSize(c.Name, 32, "name.name").Message("名字程度必须在%d-%d之内", 5, 32)

	if _, ok := v.ErrorsMap["name"]; !ok && models.DefaultCloudPlatFormManager.GetByName(c.Name) != nil {
		v.SetError("name", "名称已存在")
	}

	v.MinSize(c.Addr, 1, "name.name").Message("地址不能为空,且长程度必须在%d-%d之内", 1, 1024)
	v.MaxSize(c.Addr, 1024, "name.name").Message("地址不能为空,且长程度必须在%d-%d之内", 1, 1024)

	v.MinSize(c.Region, 1, "Region.Region").Message("Region不能为空,且长程度必须在%d-%d之内", 1, 64)
	v.MaxSize(c.Region, 64, "Region.Region").Message("Region不能为空,且长程度必须在%d-%d之内", 1, 64)

	v.MinSize(c.AccessKey, 1, "access_key.access_key").Message("access_key,且长程度必须在%d-%d之内", 1, 1024)
	v.MaxSize(c.AccessKey, 1024, "access_key.access_key").Message("access_key,且长程度必须在%d-%d之内", 1, 1024)

	v.MinSize(c.SecretKey, 1, "secret_key.nasecret_key").Message("secret_key不能为空,且长程度必须在%d-%d之内", 1, 1024)
	v.MaxSize(c.SecretKey, 1024, "secret_key.secret_key").Message("secret_key不能为空,且长程度必须在%d-%d之内", 1, 1024)

	v.MaxSize(c.Remark, 1024, "remark.remark").Message("remark长程度必须在%d之内", 1024)

	//验证类型
	if sdk, ok := cloud.DefaultManager.Cloud(c.Type); !ok {
		v.SetError("type", "类型错误")
	} else if !v.HasErrors() {
		sdk.Init(c.Addr, c.Region, c.AccessKey, c.SecretKey)
		if sdk.TestConnect() != nil {
			v.SetError("type", "配置参数错误")
		}
	}
}
