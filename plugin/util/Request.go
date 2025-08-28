package util

/*
@name 网络请求，数据爬取
*/
import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
网络请求 数据爬取
*/

type RequestInfo struct {
	Uri    string      `json:"uri"`    //请求url地址
	Params url.Values  `json:"param"`  //请求参数
	Header http.Header `json:"header"` //请求头部数据
	Resp   []byte      `json:"resp"`   //响应结果数据
	Err    string      `json:"err"`    //错误信息
}

// 爬虫实例
var Client = CreateClient()

// 记录上次的请求url
var ReferUrl string

// 初始化请求客户端
func CreateClient() *colly.Collector {
	c := colly.NewCollector()
	//访问深度
	c.MaxDepth = 1
	//可重复访问
	c.AllowURLRevisit = true
	//设置超时时间
	c.SetRequestTimeout(20 * time.Second)
	//发起请求之前调用的方法
	c.OnRequest(func(request *colly.Request) {
		//伪造请求头部
		request.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		//之前的url没有或者主域名不一致
		if len(ReferUrl) <= 0 || strings.Contains(ReferUrl, request.URL.Host) {
			ReferUrl = ""
		}
		request.Headers.Set("Referer", ReferUrl)
	})
	//异常处理
	c.OnError(func(response *colly.Response, err error) {
		log.Printf("爬虫抓取请求错误，URL： %s Error:%s\n", response.Request.URL, err)
	})
	return c

}

// 请求数据方法
func ApiGet(r *RequestInfo) {
	visit_url := fmt.Sprintf("%s?%s", r.Uri, r.Params.Encode())
	resp, err := CacheGet(visit_url)
	if err != nil {
		//正常从远端
		ApiGetRemote(r)
		err = CacheSave(visit_url, r.Resp)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("save cache ... ", visit_url)
	} else {
		//存入请求结构
		r.Resp = resp
		fmt.Println("get cache ... ", visit_url)
	}
}
func ApiGetRemote(r *RequestInfo) {
	if r.Header != nil {
		t, err := strconv.Atoi(r.Header.Get("timeout"))
		if err != nil && t > 0 {
			Client.SetRequestTimeout(time.Duration(t) * time.Second)
		}
	}
	extensions.RandomUserAgent(Client)
	//响应钩子 把爬出来的内容body放到Resp里
	Client.OnResponse(func(response *colly.Response) {
		if (response.StatusCode == 200 || (response.StatusCode >= 300 && response.StatusCode <= 399)) && len(response.Body) > 0 {
			r.Resp = response.Body
		} else {
			r.Resp = []byte{}
		}
		//将这次的响应的url保存到上次请求url变量中，给下次用
		ReferUrl = response.Request.URL.String()
	})
	//开始访问
	visit_url := fmt.Sprintf("%s?%s", r.Uri, r.Params.Encode())
	fmt.Println("start visit website url:", visit_url)
	err := Client.Visit(visit_url)
	if err != nil {
		r.Err = err.Error()
		log.Println("爬虫失败：", err)
	}
}
