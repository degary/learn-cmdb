package k8s

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type ServicePort struct {
	Id          int            `orm:"column(id)" json:"id"`
	Name        string         `orm:"column(name)" json:"name"`
	Protocol    string         `orm:"column(protocol)" json:"protocol"`
	Port        int            `orm:"column(port)" json:"port"`
	TargetPort  int            `orm:"column(target_port)" json:"target_port"`
	NodePort    int            `orm:"column(node_port)" json:"node_port"`
	Service     *ServiceModels `orm:"rel(fk);column(service)" json:"service"`
	PortNumber  int            `orm:"column(port_number)" json:"port_number"`
	CreatedTime *time.Time     `orm:"column(created_time);auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime *time.Time     `orm:"column(updated_time);auto_now;type(datetime)" json:"updated_time"`
	DeletedTime *time.Time     `orm:"column(deleted_time);null;type(datetime)" json:"deleted_time"`
}

type ServiceModels struct {
	Id          int            `orm:"column(id)" json:"id"`
	Uid         string         `orm:"column(uid);size(256)" json:"uid"`
	Name        string         `orm:"column(name)" json:"name"`
	Namespace   string         `orm:"column(namespace)" json:"namespace"`
	Port        []*ServicePort `orm:"reverse(many)" json:"port"`
	ClusterIp   string         `orm:"column(cluster_ip)" json:"cluster_ip"`
	CreatedTime *time.Time     `orm:"column(created_time);auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime *time.Time     `orm:"column(updated_time);auto_now;type(datetime)" json:"updated_time"`
	DeletedTime *time.Time     `orm:"column(deleted_time);null;type(datetime)" json:"deleted_time"`
}

type ServiceManager struct{}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

func (c *ServiceManager) Query(q string, start int64, length int) ([]*ServiceModels, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&ServiceModels{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("namespace__icontains", q)
		query = query.Or("selector__icontains", q)
		condition = condition.AndCond(query)

		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*ServiceModels
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result) //需要增加RelatedSel() 用于返回port属性的值,否则 port
	svcPorts := []*ServicePort{}
	for _, svc := range result {
		q := ormer.QueryTable("service_port")
		q.Filter("service", svc).All(&svcPorts)
		svc.Port = svcPorts
	}
	return result, total, qtotal
}

func (c *ServiceManager) Sync(svc *corev1.Service) error {
	ormer := orm.NewOrm()
	portNumber := len(svc.Spec.Ports)
	uid := strings.Split(string(svc.UID), "-")[0]
	service := ServiceModels{
		Name: svc.Name,
		Uid:  uid,
	}

	if _, _, err := ormer.ReadOrCreate(&service, "Uid", "Name"); err != nil {
		fmt.Println(err)
		return err
	}

	service.Namespace = svc.Namespace
	service.ClusterIp = svc.Spec.ClusterIP
	ormer.Update(&service)

	for n, p := range svc.Spec.Ports {
		s := ServicePort{
			Name:    p.Name,
			Service: &service,
		}
		if _, _, err := ormer.ReadOrCreate(&s, "Name", "Service"); err != nil {
			fmt.Println(err)
			return err
		}
		if n == 0 {
			s.PortNumber = portNumber
		} else {
			s.PortNumber = 0
		}
		s.Port = int(p.Port)
		s.Protocol = string(p.Protocol)
		s.TargetPort = int(p.TargetPort.IntVal)
		s.NodePort = int(p.NodePort)
		ormer.Update(&s)
	}
	return nil
}

type ServicePortManager struct{}

func (c *ServicePortManager) Query(q string, start int64, length int) ([]*ServicePort, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&ServicePort{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("namespace__icontains", q)
		query = query.Or("selector__icontains", q)
		condition = condition.AndCond(query)

		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*ServicePort
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result) //需要增加RelatedSel() 用于返回port属性的值,否则 port
	return result, total, qtotal
}

func NewServicePortManager() *ServicePortManager {
	return &ServicePortManager{}
}
