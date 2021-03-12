package forms

type FormUser struct {
	Id       int    `orm:"column(id);" json:"id" form:"id"`
	Name     string `orm:"column(name);size(32)" json:"name" form:"name"`
	Gender   int    `orm:"column(gender);default(0)" json:"gender" form:"gender"`
	Tel      string `orm:"column(tel);size(1024)" json:"tel" form:"tel"`
	Birthday string `orm:"column(birthday);null;default(null)" json:"birthday" form:"birthday"`
	Email    string `orm:"column(email);size(1024);default(null)" json:"email" form:"email"`
	Addr     string `orm:"column(addr);size(1024);default(null)" json:"addr" form:"addr"`
	Remark   string `orm:"column(remark);size(1024);default(null)" json:"remark" form:"remark"`
	Status   int    `orm:"column(status);" json:"status" form:"status"`
}
