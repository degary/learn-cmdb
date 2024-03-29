package k8s

import (
	"github.com/astaxie/beego/orm"
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

type DeploymentModels struct {
	Id            int    `orm:"column(id)" json:"id"`
	Uid           string `orm:"column(uid)" json:"uid"`
	Name          string `orm:"column(name)" json:"name"`
	Namespace     string `orm:"column(namespace)" json:"namespace"`
	Image         string `orm:"column(image)" json:"image"`
	ReadyReplicas int    `orm:"column(ready_replicas)" json:"ready_replicas"`
	Replicas      int    `orm:"column(replicas)" json:"replicas"`

	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now;type(datetime)" json:"updated_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);null;type(datetime)" json:"deleted_time"`
}

type DeploymentManager struct{}

func (c *DeploymentManager) Query(q string, start int64, length int) ([]*DeploymentModels, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&DeploymentModels{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("name__icontains", q)
		query = query.Or("namespace__icontains", q)
		query = query.Or("image__icontains", q)
		condition = condition.AndCond(query)

		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*DeploymentModels
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result) //需要增加RelatedSel() 用于返回User属性的值,否则 User属性值为空
	return result, total, qtotal
}

func (c *DeploymentManager) Sync(deploy *appsv1.Deployment) error {
	deploys := &DeploymentModels{Uid: string(deploy.ObjectMeta.UID)}
	ormer := orm.NewOrm()
	if _, _, err := ormer.ReadOrCreate(deploys, "Uid"); err != nil {
		return err
	}

	deploys.Uid = string(deploy.ObjectMeta.UID)
	deploys.Name = deploy.Name
	deploys.Namespace = deploy.Namespace
	deploys.Replicas = int(*deploy.Spec.Replicas)
	deploys.ReadyReplicas = int(deploy.Status.ReadyReplicas)
	deploys.Image = deploy.Spec.Template.Spec.Containers[0].Image

	_, err := ormer.Update(deploys)
	return err
}

func (c *DeploymentManager) GetById(id int) (*DeploymentModels, error) {
	ormer := orm.NewOrm()
	deploy := DeploymentModels{Id: id}
	err := ormer.Read(&deploy)
	if err != nil {
		return nil, err
	}
	return &deploy, nil
}

func (c *DeploymentManager) SyncDeployStatus(t time.Time, namespace string) {
	ormer := orm.NewOrm()
	//由于updatedtime字段 是 auto_now类型,在操作时,会自动更新此字段,所以 通过比较 updatetime字段的值比当前时间小 来判断此服务器是否被删除
	ormer.QueryTable(&DeploymentModels{}).Filter("Namespace__exact", namespace).Filter("DeletedTime__isnull", true).Filter("UpdatedTime__lt", t).Update(orm.Params{"DeletedTime": t})
	ormer.QueryTable(&DeploymentModels{}).Filter("Namespace__exact", namespace).Filter("DeletedTime__isnull", false).Filter("UpdatedTime__gte", t).Update(orm.Params{"DeletedTime": nil})
}

func NewDeploymentManager() *DeploymentManager {
	return &DeploymentManager{}
}

func init() {
	orm.RegisterModel(&DeploymentModels{}, &ServiceModels{}, &ServicePort{})
}
