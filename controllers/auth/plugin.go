package auth

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"github.com/degary/learn-cmdb/forms"
	"github.com/degary/learn-cmdb/models"
	"net/http"
	"strings"
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
	if c.Ctx.Input.IsAjax() {
		//ajax请求
		c.Data["json"] = map[string]interface{}{
			"code":   401,
			"text":   "请进行登录",
			"result": nil,
		}
		c.ServeJSON()

	} else {
		//http请求
		c.Redirect(beego.URLFor(beego.AppConfig.String("login")), http.StatusFound)
	}

	//ajax请求
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
				c.Redirect(beego.URLFor(beego.AppConfig.String("home")), http.StatusFound)
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

//token
type Token struct{}

func (t *Token) Is(ctx *context.Context) bool {
	//如果在请求头里,没有Authentication,则认为使用的是session
	return strings.ToLower(strings.TrimSpace(ctx.Input.Header("Authentication"))) == "token"
}
func (t *Token) Name() string {
	return "token"
}

func (t *Token) IsLogin(c *LoginRequiredController) *models.User {
	accesskey := strings.TrimSpace(c.Ctx.Input.Header("AccessKey"))
	secretkey := strings.TrimSpace(c.Ctx.Input.Header("SecretKey"))

	if token := models.DefaultTokenManager.GetByKey(accesskey, secretkey); token != nil && token.User.DeletedTime == nil {
		return token.User
	}
	return nil
}
func (t *Token) GoToLoginPage(c *LoginRequiredController) {
	c.Data["json"] = map[string]interface{}{
		"code":   403,
		"text":   "Please use a currect token",
		"result": nil,
	}
	c.ServeJSON()
}
func (t *Token) Login(c *AuthController) bool {
	c.Data["json"] = map[string]interface{}{
		"code":   403,
		"text":   "请使用token请求API",
		"result": nil,
	}
	c.ServeJSON()
	return false
}
func (t *Token) Logout(c *AuthController) {
	c.Data["json"] = map[string]interface{}{
		"code":   200,
		"text":   "退出登录成功",
		"result": nil,
	}
	c.ServeJSON()
}

func init() {
	DefaultManager.Register(new(Session))
	DefaultManager.Register(new(Token))
}
