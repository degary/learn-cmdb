package auth

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
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
			user := &models.User{Id: uid}
			ormer := orm.NewOrm()
			if err := ormer.Read(user); err == nil {
				return user
			}
		}
	}
	return nil
}
func (s *Session) GoToLoginPage(c *LoginRequiredController) {
	c.Redirect("/auth/login", http.StatusFound)
}
func (s *Session) Login(c *AuthController) bool {
	if c.Ctx.Input.IsPost() {
		//表单验证
		name := c.GetString("name")
		password := c.GetString("password")
		fmt.Println(name, password)
	}

	c.TplName = "auth/login.html"
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
