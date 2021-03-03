# learn-cmdb
## 用户验证方式
    1.用户名/密码验证
    2. Token验证
        在token表中增加对应的信息
        请求时,在Header中增加 "Authentication: token","AccessKey: xxx","SecretKey: xxx"
## 使用方法:
    配置conf/db.conf 更改数据库连接地址及账号密码
    
    启动方式:
        1. 执行 go run main.go --init -v (初始化数据库)
        2. 执行 go run main.go (启动),访问IP:Port/auth/login
