package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/degary/learn-cmdb/cloud"
	_ "github.com/degary/learn-cmdb/cloud/plugins"
	"github.com/degary/learn-cmdb/models"
	"github.com/degary/learn-cmdb/models/k8s"
	_ "github.com/degary/learn-cmdb/routers"
	k8sservice "github.com/degary/learn-cmdb/services/k8s"
	"github.com/degary/learn-cmdb/utils"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

func main() {
	//初始化命令行参数
	h := flag.Bool("h", false, "help")
	help := flag.Bool("help", false, "help")
	init := flag.Bool("init", false, "init server")
	syncdb := flag.Bool("syncdb", false, "sync db")
	force := flag.Bool("force", false, "force")
	verbose := flag.Bool("v", false, "verbose")
	flag.Usage = func() {
		fmt.Println("usage: web -h")
		flag.PrintDefaults()
	}
	//解析命令行参数
	flag.Parse()
	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	//初始化日志
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/web.log"}`)
	logs.SetLevel(logs.LevelDebug)
	logs.Debug("this is a debug log")

	if !*verbose {
		logs.GetBeeLogger().DelLogger("console")
	} else {
		orm.Debug = true
	}

	//初始化orm
	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//注册数据库
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dsn"))
	//测试数据库连接是否正常
	if db, err := orm.GetDB(); err != nil || db.Ping() != nil {
		logs.Error(err)
		os.Exit(-1)
	}
	//同步数据库
	//根据参数 选择执行流程

	switch {
	case *init:
		orm.RunSyncdb("default", *force, *verbose)
		logs.Informational("初始化数据库")
		ormer := orm.NewOrm()
		admin := models.User{Name: "admin", IsSuperuser: true}
		//Read查询admin用户是否存在,如果有报错,且报错类型为 ErrNoRaws 则表示此用户不存在
		if err := ormer.Read(&admin, "Name"); err == orm.ErrNoRows {
			//设置admin的密码
			password := utils.RandString(8)
			admin.SetPassword(password)
			if _, err = ormer.Insert(&admin); err != nil {
				logs.Error("初始化admin失败: ", err)
			} else {
				logs.Informational("初始化admin成功,默认密码: %s", password)
			}
		} else {
			logs.Informational("admin用户已存在")
		}
	case *syncdb:
		orm.RunSyncdb("default", *force, *verbose)
		logs.Informational("同步数据库")
	default:
		//同步云服务器信息
		go func() {
			for now := range time.Tick(5 * time.Second) {
				//更新云平台
				platforms, _, _ := models.DefaultCloudPlatFormManager.Query("", 0, 0)
				for _, platform := range platforms {
					if platform.IsEnable() {
						sdk, ok := cloud.DefaultManager.Cloud(platform.Type)
						if !ok {
							fmt.Println("云平台未注册")
						} else {
							sdk.Init(platform.Addr, platform.Region, platform.AccessKey, platform.SecretKey)
							if err := sdk.TestConnect(); err != nil {
								logs.Error("测试连接失败", err)
								models.DefaultCloudPlatFormManager.SyncInfo(platform, now, fmt.Sprintf("测试连接失败: %s", err.Error()))
							} else {
								instances := sdk.GetInstance()
								for _, instance := range instances {
									models.DefaultVirtualMachineManager.SyncInstance(instance, platform)
								}
								models.DefaultVirtualMachineManager.SyncInstanceStatus(now, platform)
								models.DefaultCloudPlatFormManager.SyncInfo(platform, now, "")
							}
						}
					}

				}
			}
		}()
		//同步k8s信息
		go func() {
			deploy := k8s.NewDeploymentManager()
			svc := k8s.NewServiceManager()
			client := k8sservice.NewClient("conf/config")
			for now := range time.Tick(5 * time.Second) {
				deployments, err := client.Deployments()
				if err != nil {
					fmt.Println(err)
					return
				}
				//同步deploy信息
				for _, deployment := range deployments {
					deploy.Sync(deployment)
					deploy.SyncDeployStatus(now, deployment.Namespace)
				}

				services, err := client.Services()
				if err != nil {
					fmt.Println(err)
					return
				}
				//同步service信息
				for _, service := range services {
					svc.Sync(service)
					//fmt.Printf("%#v\n",svcMo)
					svc.SyncServiceStatus(now, service.Namespace)
				}
				//同步servicePort信息
				svcs := []*k8s.ServiceModels{}
				orm.NewOrm().QueryTable(&k8s.ServiceModels{}).Filter("DeletedTime__isnull", false).All(&svcs)
				for _, svcMol := range svcs {
					svc.SyncServicePortStatus(now, svcMol)
				}
			}
		}()
		beego.Run()
	}
}
