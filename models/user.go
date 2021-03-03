package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/degary/learn-cmdb/utils"
	"time"
)

type User struct {
	Id          int        `orm:"column(id);" json:"id"`
	Name        string     `orm:"column(name);size(32)" json:"name"`
	Password    string     `orm:"column(password);size(1024);" json:"password"`
	Gender      int        `orm:"column(gender);default(0)" json:"gender"`
	Tel         string     `orm:"column(tel);size(1024)" json:"tel"`
	Birthday    *time.Time `orm:"column(birthday);null;default(null)" json:"birthday"`
	Email       string     `orm:"column(email);size(1024);default(null)" json:"email"`
	Addr        string     `orm:"column(addr);size(1024);default(null)" json:"addr"`
	Remark      string     `orm:"column(remark);size(1024);default(null)" json:"remark"`
	IsSuperuser bool       `orm:"column(is_superuser);default(false)" json:"is_superuser"`
	Status      int        `orm:"column(status);" json:"status"`
	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now" json:"updated_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);null;default(null)" json:"deleted_time"`
	//token和user是一对一的关系
	Token *Token `orm:"reverse(one)"`
}

func (u *User) SetPassword(password string) {
	u.Password = utils.Md5Salt(password, "")
}

//判断密码是否正确
func (u *User) ValidatePassword(password string) bool {
	salt, _ := utils.SplitMd5Salt(u.Password)
	return utils.Md5Salt(password, salt) == u.Password
}

//判断用户是否被锁住,锁住返回true
func (u *User) IsLocked() bool {
	return u.Status == StatusLock
}

type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (c *UserManager) GetByID(id int) *User {
	user := &User{}
	if err := orm.NewOrm().QueryTable(user).Filter("id__exact", id).Filter("DeletedTime__isnull", true).One(user); err == nil {
		return user
	}
	return nil
}

func (c *UserManager) GetByName(name string) *User {
	user := &User{}
	if err := orm.NewOrm().QueryTable(user).Filter("name__exact", name).Filter("DeletedTime__isnull", true).One(user); err == nil {
		return user
	} else {
		fmt.Println(err)
		return nil
	}

}

type Token struct {
	Id          int        `orm:"column(id)"`
	User        *User      `orm:"column(user);rel(one)"`
	AccessKey   string     `orm:"column(access_key);size(1024)"`
	SecretKey   string     `orm:"column(secret_key);size(1024)"`
	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now" json:"updated_time"`
}

type TokenManager struct{}

func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

func (t *TokenManager) GetByKey(accesskey, secretkey string) *Token {
	token := &Token{AccessKey: accesskey, SecretKey: secretkey}
	ormer := orm.NewOrm()
	if err := ormer.Read(token, "accesskey", "secretkey"); err == nil {
		ormer.LoadRelated(token, "User")
		return token
	}
	return nil
}

var DefaultUserManager = NewUserManager()
var DefaultTokenManager = NewTokenManager()

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Token))
}
