# learn-cmdb
## 用户验证方式
    1.用户名/密码验证
    2. Token验证
        在token表中增加对应的信息
        请求时,在Header中增加 "Authentication: token","AccessKey: xxx","SecretKey: xxx"
## 使用方法:
    配置conf/db.conf 更改数据库连接地址及账号密码
    注: 此程序session用的是本地file模式,可更改配置 改成redis或其他模式,详情请查看beego官网
    此程序只集成了阿里云和腾讯云,如果需要使用其他云平台,需要自己写.
        1.需要完成cloud.ICloud接口
        2.注册接口cloud.DefaultManager.Register(new(xxxCloud))
        3.在plugins/init.go中引用
    
    启动方式:
        1. 执行 go run main.go --init -v (初始化数据库)
            输出信息中会有默认的admin用户密码,请记住~
        2. 执行 go run main.go (启动),访问IP:Port/auth/login
        3. 执行 go run cloud.go 启动同步虚拟机线程
