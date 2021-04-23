package controllers

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/controllers/auth"
	"github.com/degary/learn-cmdb/forms"
	"github.com/degary/learn-cmdb/models"
	"github.com/degary/learn-cmdb/utils"
	"strings"
)

type UserPageController struct {
	LayoutController
}

func (c *UserPageController) Index() {
	c.Data["menu"] = "user_management"
	c.Data["expand"] = "system_management"
	c.LayoutSections["LayoutScript"] = "user_page/index.script.html"
	c.TplName = "user_page/index.html"
}

type UserController struct {
	auth.LoginRequiredController
}

func (c *UserController) List() {
	draw, _ := c.GetInt("draw")
	start, _ := c.GetInt64("start")
	length, _ := c.GetInt("length")
	q := strings.TrimSpace(c.GetString("q"))

	//[]*User,total,queryTotal
	users, total, queryTotal := models.DefaultUserManager.Query(q, start, length)
	c.Data["json"] = map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"result":          users,
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": queryTotal,
	}
	c.ServeJSON()
}

func (c *UserController) Create() {
	if c.Ctx.Input.IsPost() {
		json := map[string]interface{}{
			"code": 400,
			"text": "提交数据错误",
		}

		form := &forms.FormUserCreate{}
		valid := &validation.Validation{}
		if err := c.ParseForm(form); err == nil {
			if ok, err := valid.Valid(form); err != nil {
				valid.SetError("error", err.Error())
				json["result"] = valid.Errors
			} else if ok {
				user, err := models.DefaultUserManager.Create(form.Name, form.Password, form.Gender, form.BirthdayTime, form.Tel, form.Email, form.Addr, form.Remark)
				if err == nil {
					json = map[string]interface{}{
						"code":   200,
						"text":   "创建成功",
						"result": user,
					}
				} else {
					json = map[string]interface{}{
						"code": 500,
						"text": "服务器错误",
					}
				}
			} else {
				json["result"] = valid.Errors
			}
		} else {
			valid.SetError("error", err.Error())
			json["result"] = valid.Errors
		}
		c.Data["json"] = json
		c.ServeJSON()
	} else {
		//get
		c.TplName = "user/create.html"
	}
}

func (c *UserController) Modify() {
	if c.Ctx.Input.IsPost() {
		formuser := &forms.FormUser{}

		if err := c.ParseForm(formuser); err != nil {
			c.Data["json"] = map[string]interface{}{
				"code":   400,
				"text":   "编辑失败",
				"result": err,
			}
		} else {
			fmt.Printf("===================================%#v\n", formuser)
			user := &models.User{
				Id:       formuser.Id,
				Name:     formuser.Name,
				Addr:     formuser.Addr,
				Tel:      formuser.Tel,
				Birthday: utils.StrToTime(formuser.Birthday),
				Remark:   formuser.Remark,
				Email:    formuser.Email,
				Status:   formuser.Status,
				Gender:   formuser.Gender,
			}
			ormer := orm.NewOrm()
			ormer.Update(user, "name", "gender", "birthday", "tel", "email", "addr", "remark", "status")
			c.Data["json"] = map[string]interface{}{
				"code":   200,
				"text":   "创建成功",
				"result": nil,
			}
		}

		c.ServeJSON()
	} else {
		pk, _ := c.GetInt("pk")
		user := models.DefaultUserManager.GetByID(pk)
		c.Data["object"] = user
		c.TplName = "user/modify.html"
	}
}

func (c *UserController) Delete() {
	pk, _ := c.GetInt("pk")
	err := models.DefaultUserManager.DeleteById(pk)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":   400,
			"text":   err.Error(),
			"result": nil,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"code":   200,
			"text":   "删除成功",
			"result": nil,
		}
	}
	c.ServeJSON()
}

func (c *UserController) Lock() {
	pk, _ := c.GetInt("pk")
	if err := models.DefaultUserManager.SetStatusById(pk, 1); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":   400,
			"text":   err.Error(),
			"result": nil,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"code":   200,
			"text":   "锁定成功",
			"result": nil, //可以返回锁定的用户
		}
	}
	c.ServeJSON()
}

func (c *UserController) UnLock() {
	pk, _ := c.GetInt("pk")
	if err := models.DefaultUserManager.SetStatusById(pk, 0); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":   400,
			"text":   err.Error(),
			"result": nil,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"code":   200,
			"text":   "解锁成功",
			"result": nil, //可以返回锁定的用户
		}
	}
	c.ServeJSON()
}

func (c *UserController) Password() {
	if c.Ctx.Input.IsPost() {
		json := map[string]interface{}{
			"code": 400,
			"text": "提交数据错误",
		}

		form := &forms.UserPasswordForm{User: c.User}
		valid := &validation.Validation{}
		if err := c.ParseForm(form); err == nil {
			if ok, err := valid.Valid(form); err != nil {
				valid.SetError("error", err.Error())
				json["result"] = valid.Errors
			} else if ok {
				err := models.DefaultUserManager.UpdatePassword(c.User.Id, form.Password)
				if err == nil {
					json = map[string]interface{}{
						"code": 200,
						"text": "修改密码成功",
					}
				} else {
					json = map[string]interface{}{
						"code": 500,
						"text": "服务器错误",
					}
				}
			} else {
				json["result"] = valid.Errors
			}
		} else {
			valid.SetError("error", err.Error())
			json["result"] = valid.Errors
		}

		c.Data["json"] = json
		c.ServeJSON()
	} else {
		c.TplName = "user/password.html"
	}
}

type TokenController struct {
	auth.LoginRequiredController
}

func (c *TokenController) Generate() {
	if c.Ctx.Input.IsPost() {
		pk, _ := c.GetInt("pk")
		models.DefaultTokenManager.GenerateByUser(models.DefaultUserManager.GetByID(pk))
		c.Data["json"] = map[string]interface{}{
			"code":   200,
			"text":   "生成token成功",
			"result": nil, //可以返回Token
		}
		c.ServeJSON()
	} else {
		pk, _ := c.GetInt("pk")
		c.Data["object"] = models.DefaultUserManager.GetByID(pk)
		c.TplName = "token/index.html"

	}
}
