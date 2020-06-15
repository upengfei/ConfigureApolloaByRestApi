package main

import (
	"github.com/urfave/cli/v2"
	apollo "glodon.com/apollo/apolloClient"
	u "glodon.com/apollo/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	var (
		app *cli.App
	)
	app = cli.NewApp()
	app.Name = "aptool"
	app.Usage = "A tool to modify apollo configuration "
	app.Version = "1.0.0"

	app.Commands = []*cli.Command{
		{
			Name:    "modifyConf",
			Aliases: []string{"ms"},
			Usage:   "修改单个配置项命令",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "appid",
					Usage:   "项目的应用id",
					Aliases: []string{"a"},
				},
				&cli.StringFlag{
					Name:    "namespace",
					Usage:   "项目下的名称空间",
					Aliases: []string{"ns"},
					Value:   "application",
				},
				&cli.StringFlag{
					Name:    "cluster",
					Usage:   "集群名称，默认default",
					Aliases: []string{"c"},
					Value:   "default",
				},
				&cli.StringFlag{
					Name:    " section",
					Usage:   "apollo的部署环境，对应conf.ini文件的section。",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "newValue",
					Usage:   "需要修改的key的值,格式：-nv key=newvalue",
					Aliases: []string{"nv"},
				},
				&cli.StringFlag{
					Name:    "env",
					Usage:   "apollo项目里对应的环境,默认为:DEV",
					Aliases: []string{"e"},
					Value:   "DEV",
				},
				&cli.StringFlag{
					Name: "publishTitle",
					Usage: "信息发布的标题",
					Aliases: []string{"pt"},
				},
				&cli.StringFlag{
					Name: "publishContent",
					Usage: "信息发布的内容",
					Aliases: []string{"pc"},
					Value: "update",
				},
			},
			Action: func(c *cli.Context) error {
				a:=new(apollo.Apollo)
				baseData := getConfInfo(c.String("s"))
				//fmt.Println(baseData)
				apollo.InitApollo(baseData[0], baseData[1], baseData[2], a)

				a.AppId=c.String("a")
				a.NameSpace=c.String("ns")
				a.Env=c.String("e")
				a.Cluster=c.String("c")

				kv := strings.Split(c.String("nv"), "=")

				err := a.ModifyValue(kv[0], kv[1])
				if err != nil {
					return err
				}
				_,err=a.PublishMessage(c.String("pt"),c.String("pc"))
				if err!=nil{
					return err
				}


				return nil
			},
		},
		{
			Name:    "BatchEdit",
			Usage:   "通过读取csv文件方式，批量修改配置信息",
			Aliases: []string{"be"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "csv",
					Aliases: []string{"f"},
					Usage:   "csv文件路径",
				},
				&cli.StringFlag{
					Name:    "section",
					Usage:   "apollo的部署环境，对应conf.ini文件的section。",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name: "publishTitle",
					Usage: "信息发布的标题",
					Aliases: []string{"pt"},
				},
				&cli.StringFlag{
					Name: "publishContent",
					Usage: "信息发布的内容",
					Aliases: []string{"pc"},
					Value: "update",
				},
			},
			Action: func(c *cli.Context) error {
				var one sync.Once
				type apData struct {
					appId string
					nameSpace string
					env string
					cluster string
				}
				aMap:=make(map[apData]bool)
				a := new(apollo.Apollo)
				csv := u.CsvReader{
					CsvFilePath: c.String("f"),
				}

				baseData := getConfInfo(c.String("s"))
				for data := range csv.ReadContent() {
					one.Do(func() { apollo.InitApollo(baseData[0], baseData[1], baseData[2], a) })
					a.AppId = data[1]
					a.NameSpace = data[2]
					a.Env = data[0]
					a.Cluster=data[3]

					ad:=apData{
						a.AppId,
						a.NameSpace,
						a.Env,
						a.Cluster,
					}
					aMap[ad] = true
					err := a.ModifyValue(data[4], data[5])
					if err != nil {
						return err
					}
				}

				for k,_:=range aMap{

					a.AppId=k.appId
					a.NameSpace=k.nameSpace
					a.Env=k.env
					a.Cluster=k.cluster
					_,err:=a.PublishMessage(c.String("pt"),c.String("pc"))
					if err!=nil{
						return err
					}
					//}

				}

				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

}

func getConfInfo(section string) []string {
	data := make([]string, 0)
	curDir, _ := os.Getwd()
	ir := u.NewFile(filepath.Join(curDir, "conf.ini"))
	name := ir.GetValue(section, "login_user").String()
	data = append(data, name)
	password := ir.GetValue(section, "login_password").String()
	data = append(data, password)
	url := ir.GetValue(section, "login_url").String()
	data = append(data, url)
	return data

}
