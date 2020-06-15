#### 概要说明
___
aptool 是一个用来单个或批量修改配置信息的工具
> 1、源码包里同时也包含了创建、删除信息的方法，以及创建名称空间的方法，后续可根据需要进行扩展

> 2、conf.ini 文件需要与工具在同一级目录，csv不做要求，但传入时需为绝对路径。

>3、conf.ini文件里的参数是用来登录apollo服务的，根据不同的apollo部署环境，可以设置多个section，但每个section下的属性名称需保持一致，均为login_user,login_password。login_url   

#### 使用手册
___
- 通过-h命令查看子命令说明如下：
```
> aptool-win64.exe -h
NAME:
   aptool - A tool to modify apollo configuration

USAGE:
   aptool-win64.exe [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   modifyConf, ms  修改单个配置项命令
   BatchEdit, be   通过读取csv文件方式，批量修改配置信息
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```
-aptool 子命令帮助文档：
```
> aptool-win64.exe ms -h
NAME:
   aptool-win64.exe modifyConf - 修改单个配置项命令

USAGE:
   aptool-win64.exe modifyConf [command options] [arguments...]

OPTIONS:
   --appid value, -a value             项目的应用id
   --namespace value, --ns value       项目下的名称空间 (default: "application")
   --cluster value, -c value           集群名称，默认default (default: "default")
   -s value                            apollo的部署环境，对应conf.ini文件的section。
   --newValue value, --nv value        需要修改的key的值,格式：-nv key=newvalue
   --env value, -e value               apollo项目里对应的环境,默认为:DEV (default: "DEV")
   --publishTitle value, --pt value    信息发布的标题
   --publishContent value, --pc value  信息发布的内容 (default: "update")
   --help, -h                          show help (default: false)


```
```
> aptool-win64.exe be -h
NAME:
   aptool-win64.exe BatchEdit - 通过读取csv文件方式，批量修改配置信息

USAGE:
   aptool-win64.exe BatchEdit [command options] [arguments...]

OPTIONS:
   --csv value, -f value               csv文件路径
   --section value, -s value           apollo的部署环境，对应conf.ini文件的section。
   --publishTitle value, --pt value    信息发布的标题
   --publishContent value, --pc value  信息发布的内容 (default: "update")
   --help, -h                          show help (default: false)


```
#### 使用示例
- 单个配置信息修改
```
> aptool-win64.exe ms -a app123 -s apollo -nv 123=67890 -pt 测试01
```
- 批量修改
```
> aptool-win64.exe be -f ./1.csv -s apollo -pt 测试 -pc hello,world
```

#### conf文件以及csv模板
- 模板.csv
- conf.ini