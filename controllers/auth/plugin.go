package auth

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/forms"
	"github.com/degary/learn-cmdb/models"
	"net/http"
)

type Session struct {
}

func (s *Session) Is(ctx *context.Context) bool {
	//如果在请求头里,没有Authentication,则认为使用的是session
	return ctx.Input.Header("Authentication") == ""
}
func (s *Session) Name() string {
	return "session"
}

func (s *Session) IsLogin(c *LoginRequiredController) *models.User {
	if session := c.GetSession("user"); session != nil {
		if uid, ok := session.(int); ok {
			return models.DefaultUserManager.GetByID(uid)
		}
	}
	return nil
}
func (s *Session) GoToLoginPage(c *LoginRequiredController) {
	c.Redirect("/auth/login", http.StatusFound)
}
func (s *Session) Login(c *AuthController) bool {
	form := &forms.LoginForm{}
	valid := &validation.Validation{}
	if c.Ctx.Input.IsPost() {
		//表单验证
		if err := c.ParseForm(form); err != nil {
			valid.SetError("error", err.Error())
		} else {
			if ok, err := valid.Valid(form); err != nil {
				valid.SetError("error", err.Error())
			} else if ok {
				c.SetSession("user", form.User.Id)
				c.Redirect("/test/test", http.StatusFound)
				return true
			}
		}
	}

	c.TplName = "auth/login.html"
	c.Data["form"] = form
	c.Data["valid"] = valid
	return false
}
func (s *Session) Logout(c *AuthController) {
	c.DestroySession()
	c.Redirect("/auth/login", http.StatusFound)
}

type Token struct {
}

func init() {
	DefaultManager.Register(new(Session))
}
