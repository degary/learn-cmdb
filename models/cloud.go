package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/degary/learn-cmdb/cloud"
	"strings"
	"time"
)

type CloudPlatform struct {
	Id          int        `orm:"column(id)" json:"id"`
	Name        string     `orm:"column(name);size(64)" json:"name"`
	Type        string     `orm:"column(type);size(32)" json:"type"`
	Addr        string     `orm:"column(addr);size(1024)" json:"addr"`
	AccessKey   string     `orm:"column(access_key);size(1024)" json:"-"`
	SecretKey   string     `orm:"column(secret_key);size(1024)" json:"-"`
	Region      string     `orm:"column(region);size(64)" json:"region"`
	Remark      string     `orm:"column(remark);size(1024);null" json:"remark"`
	CreatedTime *time.Time `orm:"column(created_time);type(datetime);auto_now_add;" json:"created_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);type(datetime);null;" json:"deleted_time"`
	SyncTime    *time.Time `orm:"column(sync_time);type(datetime);null" json:"sync_time"`
	User        *User      `orm:"column(user);rel(fk)" json:"user"` //fk用于外键 一对多当中的"多"
	Status      int        `orm:"column(status)" json:"status"`
	Msg         string     `orm:"column(msg);size(1024)" json:"msg"`

	VirtualMachines []*VirtualMachine `orm:"reverse(many)" json:"virtual_machines"`
}

func (c *CloudPlatform) IsEnable() bool {
	return c.Status == 0 && c.DeletedTime == nil
}

type CloudPlatformManager struct {
}

func (c *CloudPlatformManager) Query(q string, start int64, length int) ([]*CloudPlatform, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&CloudPlatform{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("addr__icontains", q)
		query = query.Or("type__icontains", q)
		query = query.Or("region__icontains", q)
		query = query.Or("remark__icontains", q)
		condition = condition.AndCond(query)

		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*CloudPlatform
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result) //需要增加RelatedSel() 用于返回User属性的值,否则 User属性值为空
	return result, total, qtotal
}

func (c *CloudPlatformManager) GetByName(name string) *CloudPlatform {
	ormer := orm.NewOrm()
	var result CloudPlatform
	err := ormer.QueryTable(&CloudPlatform{}).Filter("deleted_time__isnull", true).Filter("name__exact", name).One(&result)
	if err == nil {
		return &result
	}
	return nil
}

func (c *CloudPlatformManager) Create(name, typ, addr, region, access_key, secret_key, remark string, user *User) (*CloudPlatform, error) {
	ormer := orm.NewOrm()
	result := &CloudPlatform{
		Name:      name,
		Type:      typ,
		Addr:      addr,
		Region:    region,
		Remark:    remark,
		AccessKey: access_key,
		SecretKey: secret_key,
		User:      user,
		Status:    0,
	}
	if _, err := ormer.Insert(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CloudPlatformManager) DeleteById(id int) error {
	now := time.Now()
	DefaultVirtualMachineManager.DeleteByPlatformId(id)
	_, err := orm.NewOrm().QueryTable(&CloudPlatform{}).Filter("Id__exact", id).Update(orm.Params{"DeletedTime": &now})
	return err
}

func (c *CloudPlatformManager) SyncInfo(platform *CloudPlatform, now time.Time, msg string) error {
	platform.SyncTime = &now
	platform.Msg = msg
	_, err := orm.NewOrm().Update(platform)
	return err
}

func NewCloudPlatFormManager() *CloudPlatformManager {
	return &CloudPlatformManager{}
}

type VirtualMachine struct {
	Id            int            `orm:"column(id)" json:"id"`
	Platform      *CloudPlatform `orm:"column(platform);rel(fk);" json:"platform"`
	UUID          string         `orm:"column(uuid);size(128)" json:"uuid"`
	Name          string         `orm:"column(name);size(64)" json:"name"`
	CPU           int            `orm:"column(cpu);" json:"cpu"`
	Mem           int64          `orm:"column(mem);" json:"mem"`
	Status        string         `orm:"column(status);" json:"status"`
	OS            string         `orm:"column(os);size(128)" json:"os"`
	PrivateAddrs  string         `orm:"colume(private_addrs);size(1024)" json:"private_addrs"`
	PublicAddrs   string         `orm:"column(public_addrs);size(1024)" json:"public_addrs"`
	VmCreatedTime string         `orm:"column(vm_created_time)" json:"vm_created_time"`
	VmExpiredTime string         `orm:"column(vm_expired_time)" json:"vm_expired_time"`

	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now;type(datetime)" json:"updated_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);null;type(datetime)" json:"deleted_time"`
}

type VirtualMachineManager struct{}

func NewVirtualMachineManager() *VirtualMachineManager {
	return &VirtualMachineManager{}
}

func (c *VirtualMachineManager) Query(q string, start int64, length int) ([]*VirtualMachine, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&VirtualMachine{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("platform__icontains", q)
		query = query.Or("os__icontains", q)
		query = query.Or("private_addrs__icontains", q)
		query = query.Or("public_addrs__icontains", q)
		condition = condition.AndCond(query)

		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*VirtualMachine
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result) //需要增加RelatedSel() 用于返回User属性的值,否则 User属性值为空
	return result, total, qtotal
}

func (c *VirtualMachineManager) SyncInstance(instance *cloud.Instance, platform *CloudPlatform) {
	vm := &VirtualMachine{UUID: instance.UUID, Platform: platform}
	ormer := orm.NewOrm()
	if _, _, err := ormer.ReadOrCreate(vm, "UUID", "Platform"); err != nil {
		logs.Error(err)
		return
	}
	//存在
	vm.Name = instance.Name
	vm.OS = instance.Os
	vm.CPU = instance.CPU
	vm.Mem = instance.Memory
	vm.Status = instance.Status
	vm.PublicAddrs = strings.Join(instance.PublicAddrs, ",")
	vm.PrivateAddrs = strings.Join(instance.PrivateAddrs, ",")
	vm.VmCreatedTime = instance.CreatedTime
	vm.VmExpiredTime = instance.ExpiredTime
	ormer.Update(vm)
}

func (c *VirtualMachineManager) SyncInstanceStatus(t time.Time, platform *CloudPlatform) {
	ormer := orm.NewOrm()
	//由于updatedtime字段 是 auto_now类型,在操作时,会自动更新此字段,所以 通过比较 updatetime字段的值比当前时间小 来判断此服务器是否被删除
	ormer.QueryTable(&VirtualMachine{}).Filter("Platform__exact", platform).Filter("UpdatedTime__lt", t).Update(orm.Params{"DeletedTime": t})
	ormer.QueryTable(&VirtualMachine{}).Filter("Platform__exact", platform).Filter("UpdatedTime__gte", t).Update(orm.Params{"DeletedTime": nil})
}

func (c *VirtualMachineManager) GetById(pk int) *VirtualMachine {
	vm := &VirtualMachine{}
	if err := orm.NewOrm().QueryTable(&VirtualMachine{}).RelatedSel().Filter("id__exact", pk).Filter("DeletedTime__isnull", true).One(vm); err == nil {
		return vm
	} else {
		logs.Error(err)
		return nil
	}

}

func (c *VirtualMachineManager) DeleteByPlatformId(platform int) {
	orm.NewOrm().QueryTable(&VirtualMachine{}).Filter("platform__exact", platform).Update(orm.Params{"DeletedTime": time.Now()})
}

func (c *VirtualMachineManager) Count() int {
	count, _ := orm.NewOrm().QueryTable(&VirtualMachine{}).Filter("deleted_time__isnull", true).Count()
	return int(count)
}

var DefaultCloudPlatFormManager = NewCloudPlatFormManager()
var DefaultVirtualMachineManager = NewVirtualMachineManager()

func init() {
	orm.RegisterModel(&CloudPlatform{}, &VirtualMachine{}, &Alert{})
}
