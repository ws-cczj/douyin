# douyin

<!-- PROJECT SHIELDS -->
<br>

![GitHub Repo stars](https://img.shields.io/github/stars/ws-cczj/douyin?style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/ws-cczj/douyin?style=plastic)
![GitHub forks](https://img.shields.io/github/forks/ws-cczj/douyin?style=plastic)
![GitHub contributors](https://img.shields.io/github/contributors/ws-cczj/douyin)


<!-- PROJECT LOGO -->
<br />


### 上手指南

#### 开发前的配置要求

1. go 1.19.1
2. MySQL(数据库sql文件在models包中)
3. 搭建Redis、mongodb环境
4. 拥有ffmpeg可执行文件

#### 如何配置
在conf包下的yaml文件中将地址进行修改配置即可

#### 如何运行
编译文件: `go build main.go`

直接运行: `go run main.go`

通过make一步完成: `make`
#### 如何拉取该项目
```sh
git clone https://github.com/ws-cczj/douyin.git
```

### 文件目录说明

```
douyin 
├─cache 缓存
├─conf 配置
├─consts 全局常量
├─database 数据库
│  ├─models
│  └─mongodb
├─handlers 响应处理
│  ├─comment
│  ├─common
│  ├─favor
│  ├─message
│  ├─relation
│  ├─user
│  └─video
├─middleware 中间件
├─pkg 全局库
│  ├─document 敏感词库
│  ├─e 统一错误库
│  ├─logger 日志库
│  └─utils 工具
├─public 公共文件
│  └─pic
├─router 路由
└─service 服务
    ├─comment
    ├─favor
    ├─message
    ├─relation
    ├─user
    └─video
```

### 数据库表设计
![database](http://cdn.cczjblog.top/cczjBlog-img/douyin_database.png-cczjImage)

- 避免使用外键关联.外键关联会导致删除时出现连锁问题,并且会导致插入效率变慢,处理比较麻烦.
- 避免硬删除.硬删除会造成主键在B+树中不连续,造成查询效率慢.碎片化等问题.
- 因为是软删除,所以多处使用复合索引加快查询效率.

### 测试
使用`go_test`对部分代码进行测试输出

使用`pprof`配合go内置工具生成火焰图

`go tool pprof -http=:4399 http://192.168.43.219:8080/debug/pprof/profile`
![flameGraph](http://cdn.cczjblog.top/cczjBlog-img/douyin_flamegraph.png-cczjImage)
#### 使用`go-wrk` 工具进行测试
`go-wrk -t=8 -c=100 -n=20000 http://xxx.xxx.xxx.xxx:8080/douyin/xxxx/`

经过测试，所有的请求在`2W`次请求测试中均达到无错误并且总响应速度在`1.5s`左右


### 使用到的技术
框架相关:
- [Gin](https://gin-gonic.com/docs/)

数据库相关:
- [MySQL](https://dev.mysql.com/doc/)
- [sqlx](https://github.com/jmoiron/sqlx)
- [mongodb](https://www.mongodb.com/docs/drivers/go/current/)

工具相关:
- [ffmpeg](https://ffmpeg.org/documentation.html)

中间件相关:
- [go-redis](https://juejin.cn/post/7027347979065360392)
- [JWT](https://jwt.io/introduction)

### TODO
如果有机会会去尝试一下将该项目改为微服务模式