package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/degary/learn-cmdb/cloud"
	_ "github.com/degary/learn-cmdb/cloud/plugins"
	_ "github.com/degary/learn-cmdb/routers"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"

	"github.com/degary/learn-cmdb/models"
)

func main() {
	//初始化命令行参数
	h := flag.Bool("h", false, "help")
	help := flag.Bool("help", false, "help")
	verbose := flag.Bool("v", false, "verbose")
	flag.Usage = func() {
		fmt.Println("usage: cloud -h")
		flag.PrintDefaults()
	}
	//解析命令行参数
	flag.Parse()
	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	//初始化日志
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/cloud.log"}`)
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

	for now := range time.Tick(5 * time.Second) {
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

}
