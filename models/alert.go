package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Alert struct {
	Id          int        `orm:"column(id);" json:"id"`
	Fingerprint string     `orm:"column(fingerprint);size(128)"`
	Instance    string     `orm:"column(instance)" json:"instance"`
	AlertName   string     `orm:"column(alert_name);size(256)" json:"alert_name"`
	Severity    string     `orm:"column(severity);size(32)" json:"severity"`
	Status      string     `orm:"column(status)" json:"status"`
	Labels      string     `orm:"column(labels);type(longtext)" json:"labels"`
	Summary     string     `orm:"column(summary);type(text)" json:"summary"`
	Description string     `orm:"column(description);type(text)" json:"description"`
	Annotations string     `orm:"column(annotations);type(longtext)" json:"annotations"`
	StartsAt    *time.Time `orm:"column(starts_at);type(datetime)" json:"starts_at"`
	EndsAt      *time.Time `orm:"column(ends_at);type(datetime);null" json:"ends_at"`

	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now;type(datetime)" json:"updated_time"`
	DeletedTime *time.Time `orm:"column(deleted_time);null;type(datetime)" json:"deleted_time"`
}

type AlertManager struct{}

func (c *AlertManager) Query(q string, start int64, length int) ([]*Alert, int64, int64) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&Alert{})

	condition := orm.NewCondition()
	condition = condition.And("deleted_time__isnull", true)

	total, _ := qs.SetCond(condition).Count()
	qtotal := total

	if q != "" {
		query := orm.NewCondition()
		query = query.Or("status__icontains", q)
		query = query.Or("instance__icontains", q)
		query = query.Or("severity__icontains", q)
		condition = condition.AndCond(query)
		qtotal, _ = qs.SetCond(condition).Count()
	}
	var result []*Alert
	qs.SetCond(condition).RelatedSel().Limit(length).Offset(start).All(&result)
	return result, total, qtotal
}

func (c *AlertManager) Notify(alert *Alert) {
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(&Alert{})
	qs = qs.Filter("Fingerprint", alert.Fingerprint)
	qs = qs.Filter("DeletedTime__isnull", true)
	qs = qs.Filter("Status", "firing")

	if alert.Status == "firing" {
		if cnt, err := qs.Count(); err == nil && cnt == 0 {
			ormer.Insert(alert)
		}
	} else {
		qs.Update(orm.Params{
			"EndsAt": alert.EndsAt,
			"Status": alert.Status,
		})
	}
}

func newAlertManager() *AlertManager {
	return &AlertManager{}
}

var DefaultAlertManager = newAlertManager()
