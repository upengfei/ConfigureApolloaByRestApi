package apolloClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/phachon/go-logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Apollo struct {
	Burl      string //访问地址
	AppId     string // 应用id
	NameSpace string //命名空间
	Env       string //环境信息
	Cluster   string //集群名称
	Cookies   []*http.Cookie
	C         *http.Client
}

var logger *g.Logger
func init(){
	logger = g.NewLogger()

	logger.Detach("console")

	// console adapter config
	consoleConfig := &g.ConsoleConfig{
		Color:      true,  // Does the text display the color
		JsonFormat: false, // Whether or not formatted into a JSON string
		Format:     "",    // JsonFormat is false, logger message output to console format string
	}
	// add output to the console
	logger.Attach("console", g.LOGGER_LEVEL_INFO, consoleConfig)
}


//登陆apollo，并获取cookie
func InitApollo(name,password,baseUrl string,apollo *Apollo) {
	apollo.C = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	payload := make(url.Values)
	payload.Add("username", name)
	payload.Add("password", password)
	payload.Add("login-submit", "登录")

	loginUrl := fmt.Sprintf("%s/signin", baseUrl)
	req, _ := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(payload.Encode()))

	//设置header
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	req.Header["Upgrade-Insecure-Requests"] = []string{"1"}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0"}
	req.AddCookie(&http.Cookie{Name: "NG_TRANSLATE_LANG_KEY", Value: "zh-CN"})
	resp, err := apollo.C.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == 200{
		logger.Info("登录apollo环境成功!")
	}
	apollo.Cookies = resp.Cookies()
	apollo.Burl=baseUrl
	//return apollo
}

//请求头设置
func (a *Apollo) reqHeaderSetup(req *http.Request) {
	for _, cookie := range a.Cookies {
		req.AddCookie(cookie)
	}
	req.AddCookie(&http.Cookie{Name: "NG_TRANSLATE_LANG_KEY", Value: "zh-CN"})
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Set("Referer", "http://106.54.227.205/config.html")
}

//获取指定应用的命名空间下的配置的信息
func (a *Apollo) GetConfInfo() ([]map[string]interface{}, error) {
	var (
		err         error
		resp        *http.Response
		req         *http.Request
		resultSlice []map[string]interface{}
	)

	apiUrl := fmt.Sprintf("/apps/%s/envs/%s/clusters/%s/namespaces/%s", a.AppId, a.Env, a.Cluster, a.NameSpace)
	//fmt.Println(apiUrl)
	req, err = http.NewRequest(http.MethodGet, a.Burl+apiUrl, nil)
	if err != nil {
		return nil, err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	resp, err = a.C.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	//result:=make([]interface{},0)
	result := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)

	resultSlice = make([]map[string]interface{}, 0)
	for _, v := range result["items"].([]interface{}) {
		resultSlice = append(resultSlice, v.(map[string]interface{})["item"].(map[string]interface{}))
	}

	return resultSlice, nil
}

//通过key获取该配置的具体信息
func (a *Apollo) GetSpecificConfInfo(key string) (map[string]interface{}, error) {
	result, err := a.GetConfInfo()
	if err!=nil{

		panic(err)
	}
	for _, conf := range result {
		if conf["key"] == key {
			message,_:=json.Marshal(conf)
			logger.Infof("获取%s的配置信息：%s",key,string(message))
			return conf, nil
		}
	}
	return nil, errors.New("configuration is not exists")
}

//添加配置信息
func (a *Apollo) AddNewConf(key, value string, comment ...string) (map[string]interface{}, error) {
	var (
		err error

		resp *http.Response
		req  *http.Request
	)
	apiUrl := fmt.Sprintf("/apps/%s/envs/%s/clusters/%s/namespaces/%s/item", a.AppId, a.Env, a.Cluster, a.NameSpace)
	//fmt.Println(api_url)
	payload := make(map[string]interface{})
	payload["addItemBtnDisabled"] = true

	if len(comment) > 0 {
		payload["comment"] = comment[0]
	}

	payload["key"] = key
	payload["value"] = value
	payload["tableViewOperType"] = "create"
	data, _ := json.Marshal(payload)
	req, err = http.NewRequest(http.MethodPost, a.Burl+apiUrl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = a.C.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	result := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Infof("添加配置信息成功，返回信息:%s",string(body))
	_ = json.Unmarshal(body, &result)
	return result, nil
}

//修改指定配置
func (a *Apollo) ModifyValue(key string, newValue string) error {
	var (
		err  error
		info map[string]interface{}
		resp *http.Response
		req  *http.Request
	)
	info, err = a.GetSpecificConfInfo(key)
	if err != nil {
		logger.Errorf("获取配置信息错误:%s",err)
		panic(err)
	}
	info["value"] = newValue
	data, _ := json.Marshal(info)
	apiUrl := fmt.Sprintf("/apps/%s/envs/%s/clusters/%s/namespaces/%s/item", a.AppId, a.Env, a.Cluster, a.NameSpace)
	req, err = http.NewRequest(http.MethodPut, a.Burl+apiUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = a.C.Do(req)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()
	if resp.StatusCode == 200 {
		logger.Infof("修改当前key:[%s]的配置信息成功!",key)
	}

	return nil
}

//删除指定配置
func (a *Apollo) DeleteConfig(key string) error {
	var (
		err  error
		info map[string]interface{}
		resp *http.Response
		req  *http.Request
	)
	info, err = a.GetSpecificConfInfo(key)
	if err != nil {
		panic(err)
	}
	apiUrl := fmt.Sprintf("/apps/%s/envs/%s/clusters/%s/namespaces/%s/items/%d",
		a.AppId,
		a.Env,
		a.Cluster,
		a.NameSpace,
		int(info["id"].(float64)),
	)
	//fmt.Println(apiUrl)
	req, err = http.NewRequest(http.MethodDelete, a.Burl+apiUrl, nil)
	if err != nil {
		return err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = a.C.Do(req)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()
	if resp.StatusCode == 200{
		logger.Infof("删除key:[%s]的配置信息成功",key)
	}

	return nil
}

//发布指定namespace配置信息
func (a *Apollo) PublishMessage(releaseTitle, releaseContent string) (map[string]interface{}, error) {
	var (
		err  error
		data map[string]interface{}
		resp *http.Response
		req  *http.Request
	)
	data = make(map[string]interface{})
	data["isEmergencyPublish"] = false
	data["releaseComment"] = releaseContent
	data["releaseTitle"] = releaseTitle
	payload, _ := json.Marshal(data)
	apiUrl := fmt.Sprintf("/apps/%s/envs/%s/clusters/%s/namespaces/%s/releases",
		a.AppId,
		a.Env,
		a.Cluster,
		a.NameSpace,
	)
	req, err = http.NewRequest(http.MethodPost, a.Burl+apiUrl, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = a.C.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	result := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Infof("【AppId:%s】-【namespace:[%s]】发布信息成功，返回值:%s",a.AppId,a.NameSpace,string(body))
	_ = json.Unmarshal(body, &result)
	return result, nil
}

//创建namespace
func (a *Apollo) CreateNameSpace(name string, comment string, isPublic bool) (map[string]interface{}, error) {
	var (
		err    error
		data   map[string]interface{}
		resp   *http.Response
		req    *http.Request
		apiUrl string
	)
	apiUrl = fmt.Sprintf("/apps/%s/appnamespaces?appendNamespacePrefix=true", a.AppId)
	data = make(map[string]interface{})
	data["appId"] = a.AppId
	data["comment"] = comment
	data["format"] = "properties"
	data["isPublic"] = isPublic
	data["name"] = name
	payload, _ := json.Marshal(data)
	req, err = http.NewRequest(http.MethodPost, a.Burl+apiUrl, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	//设置请求头
	a.reqHeaderSetup(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = a.C.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	result := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)
	return result, nil

}

//删除namespace
//TODO
