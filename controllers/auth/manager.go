package auth

import (
	"github.com/astaxie/beego/context"
	"github.com/degary/learn-cmdb/models"
)

type AuthPlugin interface {
	Name() string
	Is(*context.Context) bool
	IsLogin(*LoginRequiredController) *models.User
	GoToLoginPage(*LoginRequiredController)
	Login(*AuthController) bool
	Logout(*AuthController)
}

type Manager struct {
	plugins map[string]AuthPlugin
}

func NewManager() *Manager {
	return &Manager{
		plugins: map[string]AuthPlugin{},
	}
}

func (m *Manager) Register(p AuthPlugin) {
	name := p.Name()
	m.plugins[name] = p
}

func (m *Manager) GetPlugin(c *context.Context) AuthPlugin {
	//遍历plugins,执行Is方法,如果为true,则返回该plugin
	for _, plugin := range m.plugins {
		if plugin.Is(c) {
			return plugin
		}
	}
	return nil
}

func (m *Manager) IsLogin(c *LoginRequiredController) *models.User {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		return plugin.IsLogin(c)
	}
	return nil
}

func (m *Manager) GoToLoginPage(c *LoginRequiredController) {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		plugin.GoToLoginPage(c)
	}
}

func (m *Manager) Login(c *AuthController) bool {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		return plugin.Login(c)
	}
	return false
}

func (m *Manager) Logout(c *AuthController) {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		plugin.Logout(c)
	}
}

var DefaultManager = NewManager()
