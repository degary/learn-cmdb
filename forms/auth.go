package forms

import (
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/models"
	"strings"
)

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	User     *models.User
}

func (c *LoginForm) Valid(v *validation.Validation) {
	c.Username = strings.TrimSpace(c.Username)
	c.Password = strings.TrimSpace(c.Password)

	if c.Username == "" || c.Password == "" {
		v.SetError("error", " 用户名或密码不能为空")
	} else {
		//通过username查找用户
		if user := models.DefaultUserManager.GetByName(c.Username); user == nil {
			v.SetError("error", "用户不存在")
		} else if !user.ValidatePassword(c.Password) {
			v.SetError("error", "用户名或密码错误")
		} else if user.IsLocked() { //检查是否为lock状态
			v.SetError("error", "用户已被锁定")
		} else {
			c.User = user
		}

	}
}
