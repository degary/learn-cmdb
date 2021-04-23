package forms

import (
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/models"
	"strings"
	"time"
)

type FormUser struct {
	Id       int    `form:"id"`
	Name     string `form:"name"`
	Gender   int    `form:"gender"`
	Tel      string `form:"tel"`
	Birthday string `form:"birthday"`
	Email    string `form:"email"`
	Addr     string `form:"addr"`
	Remark   string `form:"remark"`
	Status   int    `form:"status"`
}

type FormUserCreate struct {
	Password       string `form:"password"`
	PasswordVerify string `form:"passwordVerify"`
	Name           string `form:"name"`
	Gender         int    `form:"gender"`
	Tel            string `form:"tel"`
	Birthday       string `form:"birthday"`
	Email          string `form:"email"`
	Addr           string `form:"addr"`
	Remark         string `form:"remark"`

	BirthdayTime *time.Time
}

func (c *FormUserCreate) Valid(v *validation.Validation) {
	c.Name = strings.TrimSpace(c.Name)
	c.Password = strings.TrimSpace(c.Password)
	c.PasswordVerify = strings.TrimSpace(c.PasswordVerify)
	c.Tel = strings.TrimSpace(c.Tel)
	c.Email = strings.TrimSpace(c.Email)
	c.Remark = strings.TrimSpace(c.Remark)
	c.Addr = strings.TrimSpace(c.Addr)

	v.AlphaDash(c.Name, "name.name").Message("用户名只能由数字、英文字母、中划线和下划线组成")
	v.MinSize(c.Name, 5, "name.name").Message("用户名长度必须在%d-%d之内", 5, 32)
	v.MaxSize(c.Name, 32, "name.name").Message("用户名长度必须在%d-%d之内", 5, 32)

	if _, ok := v.ErrorsMap["name"]; !ok && models.DefaultUserManager.GetByName(c.Name) != nil {
		v.SetError("name", "用户名已存在")
	}

	v.MinSize(c.Password, 6, "password.password").Message("密码长度必须在%d-%d之内", 6, 32)
	v.MaxSize(c.Password, 32, "password.password").Message("密码长度必须在%d-%d之内", 6, 32)

	if c.PasswordVerify != c.Password {
		v.SetError("passwordVerify", "两次输入的密码不一致")
	}

	v.Range(c.Gender, 0, 1, "gender.gender").Message("性别选择不正确")

	if Birthday, err := time.Parse("2006-01-02", c.Birthday); err != nil {
		v.SetError("birthday", "出生日期不正确")
	} else {
		c.BirthdayTime = &Birthday
	}

	v.Phone(c.Tel, "tel.tel").Message("电话格式不正确")
	v.Email(c.Email, "email.email").Message("邮箱格式不正确")

	v.MaxSize(c.Addr, 512, "addr.addr").Message("住址长度必须在512个字符之内")
	v.MaxSize(c.Remark, 512, "remark.remark").Message("备注长度必须在512个字符之内")
}

type UserPasswordForm struct {
	OldPassword    string `form:"oldPassword"`
	Password       string `form:"password"`
	PasswordVerify string `form:"passwordVerify"`

	User *models.User
}

func (f *UserPasswordForm) Valid(v *validation.Validation) {
	f.OldPassword = strings.TrimSpace(f.OldPassword)
	f.Password = strings.TrimSpace(f.Password)
	f.PasswordVerify = strings.TrimSpace(f.PasswordVerify)

	if !f.User.ValidatePassword(f.OldPassword) {
		v.SetError("oldPassword", "密码不正确")
	}

	v.MinSize(f.Password, 6, "password.password").Message("密码长度必须在%d-%d之内", 6, 32)
	v.MaxSize(f.Password, 32, "password.password").Message("密码长度必须在%d-%d之内", 6, 32)

	if f.PasswordVerify != f.Password {
		v.SetError("passwordVerify", "两次输入密码不一致")
	}
}
