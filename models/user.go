package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/degary/learn-cmdb/utils"
	"time"
)

type User struct {
	Id          int        `orm:"column(id);" json:"id" form:"id"`
	Name        string     `orm:"column(name);size(32)" json:"name" form:"name"`
	Password    string     `orm:"column(password);size(1024);" json:"-"`
	Gender      int        `orm:"column(gender);default(0)" json:"gender" form:"gender"`
	Tel         string     `orm:"column(tel);size(1024)" json:"tel" form:"tel"`
	Birthday    *time.Time `orm:"column(birthday);type(date);null;default(null)" json:"birthday" form:"birthday"`
	Email       string     `orm:"column(email);size(1024);default(null)" json:"email" form:"email"`
	Addr        string     `orm:"column(addr);size(1024);default(null)" json:"addr" form:"addr"`
	Remark      string     `orm:"column(remark);size(1024);default(null)" json:"remark" form:"remark"`
	IsSuperuser bool       `orm:"column(is_superuser);default(false)" json:"is_superuser"`
	Status      int        `orm:"column(status);" json:"status" form:"status"`
	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now" json:"updated_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);null;default(null)" json:"deleted_time"`
	//token和user是一对一的关系
	Token          *Token           `orm:"reverse(one)" json:"token"`
	CloudPlatForms []*CloudPlatform `orm:"reverse(many)" json:"cloud_platforms"`
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
	ormer := orm.NewOrm()
	if err := ormer.QueryTable(user).Filter("id__exact", id).Filter("DeletedTime__isnull", true).One(user); err == nil {
		ormer.LoadRelated(user, "Token")
		return user
	}
	return nil
}

func (c *UserManager) DeleteById(id int) error {
	user := new(User)
	qs := orm.NewOrm().QueryTable(user)
	qs.Filter("id__exact", id).Update(orm.Params{"deleted_time": time.Now()})
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

func (c *UserManager) SetStatusById(id, status int) error {
	user := new(User)
	orm.NewOrm().QueryTable(user).Filter("id__exact", id).Update(orm.Params{"status": status})
	return nil
}

func (c *UserManager) Query(q string, start int64, length int) ([]*User, int64, int64) {
	ormer := orm.NewOrm()
	queryset := ormer.QueryTable(&User{})
	condition := orm.NewCondition()

	//设置条件==> deleted_time 为空(此用户没有被删除)
	condition = condition.And("deleted_time__isnull", true)

	total, _ := queryset.SetCond(condition).Count()

	qtotal := total
	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("tel__icontains", q)
		query = query.Or("addr__icontains", q)
		query = query.Or("email__icontains", q)
		query = query.Or("remark__icontains", q)
		condition = condition.AndCond(query)
		qtotal, _ = queryset.SetCond(condition).Count()
	}
	var users []*User
	queryset.RelatedSel().SetCond(condition).Limit(length).Offset(start).All(&users)
	return users, total, qtotal

}

func (c *UserManager) Create(name, password string, gender int, birthday *time.Time, tel, email, addr, remark string) (*User, error) {
	ormer := orm.NewOrm()
	user := &User{
		Name:     name,
		Gender:   gender,
		Birthday: birthday,
		Tel:      tel,
		Email:    email,
		Addr:     addr,
		Remark:   remark,
	}
	user.SetPassword(password)
	if _, err := ormer.Insert(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *UserManager) UpdatePassword(id int, password string) error {
	ormer := orm.NewOrm()
	if user := c.GetByID(id); user != nil {
		user.SetPassword(password)
		if _, err := ormer.Update(user, "Password"); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("操作对象不存在")
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

func (t *TokenManager) GenerateByUser(user *User) *Token {
	ormer := orm.NewOrm()
	token := &Token{User: user}
	//如果查询无果,则插入,反之 则更新
	if ormer.Read(token, "User") == orm.ErrNoRows {
		token.AccessKey = utils.RandString(32)
		token.SecretKey = utils.RandString(32)
		ormer.Insert(token)
	} else {
		token.AccessKey = utils.RandString(32)
		token.SecretKey = utils.RandString(32)
		ormer.Update(token)
	}
	return token
}

var DefaultUserManager = NewUserManager()
var DefaultTokenManager = NewTokenManager()

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Token))
}
