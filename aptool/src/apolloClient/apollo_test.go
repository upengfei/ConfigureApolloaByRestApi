package apolloClient

import (
	"fmt"
	"testing"
)

func TestApollo_GetConfInfo(t *testing.T) {

}

func TestApollo_GetNameSpaceInfo(t *testing.T) {
	//a:=&Apollo{
	//	AppId: "app123",
	//	NameSpace: "application",
	//	Env: "DEV",
	//}
	a:=new(Apollo)
	InitApollo("apollo","admin","http://106.54.227.205",a)
	a.AppId="app123"
	a.NameSpace="application"
	a.Env="DEV"
	a.Cluster="default"
	data,_:=a.GetConfInfo()
	fmt.Println(data)


}

func TestApollo_ModifyValue(t *testing.T) {
	a:=&Apollo{
		AppId: "app123",
		NameSpace: "application",
		Env: "DEV",
	}
	InitApollo("apollo","admin","http://106.54.227.205",a)
	_=a.ModifyValue("ceshi","jjjjj")
}

func TestApollo_DeleteConfig(t *testing.T) {
	a:=&Apollo{
		AppId: "app123",
		NameSpace: "application",
		Env: "DEV",
	}
	InitApollo("apollo","admin","http://106.54.227.205",a)
	_=a.DeleteConfig("fff")
}

func TestApollo_PublishMessage(t *testing.T) {
	a:=&Apollo{
		AppId: "app123",
		NameSpace: "application",
		Env: "DEV",
	}
	InitApollo("apollo","admin","http://106.54.227.205",a)
	data,_:=a.PublishMessage("ceshi","")
	fmt.Println(data)
}

//创建namespace
func TestApollo_CreateNameSpace(t *testing.T) {
	a:=&Apollo{
		AppId: "app123",
		NameSpace: "application",
		Env: "DEV",
	}
	InitApollo("apollo","admin","http://106.54.227.205",a)
	data,_:=a.CreateNameSpace("ceshi002","hello,world!",false)
	fmt.Println(data)
}