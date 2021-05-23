package controller

import (
	"awsSDK/app/module"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)
//
type ES struct {
	Client *elasticsearch.Client
}

func (es *ES) RequestIndex(jsonlogstr string){
	req := esapi.IndexRequest{
		Index:   module.Configs.Index,
		Body:    strings.NewReader(jsonlogstr),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		fmt.Println("request error")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("请求结果： ", res)
	defer res.Body.Close()

}

func (es *ES) SearchIndex() (Timestamp int64 ,err error) {
	//定义传入倒排请求
	searchbody :=
		`{
  "version": true,
  "size": 500,
  "sort": [
    {
      "Time": {
           "order": "desc"
      }
    }
  ]
}`
	//检索ES中的日志，指定Index 和 检索请求Body
	req := esapi.SearchRequest{
		Index: []string{module.Configs.Index},
		Body:  strings.NewReader(searchbody),
	}
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		fmt.Println("request error")
		fmt.Println(err)
		os.Exit(1)
	}
	//把检索结果中的第一个Timestamp取出，并转成int格式
	body,_ := io.ReadAll(res.Body)
	r, _ := regexp.Compile(`"Timestamp":\d{13}`)
	t:=r.FindString(string(body))
	//如果第一次执行没有Timestamp记录，则令Timestamp = 0
	if len(t) == 0 {
		return 0,err
	}
	t2 := t[12:]
	Timestamp,_ = strconv.ParseInt(t2,10,64)
	defer res.Body.Close()
	return Timestamp+1,err
}


//初始化ES客户端
var ESClient *ES
var once sync.Once

func init() {
	once.Do(
		func() {
			cfg := elasticsearch.Config{
				Addresses: []string{
					module.Configs.Esendpoint,
				},
				Username: module.Configs.Esuser,
				Password: module.Configs.Espasswd,
			}
			client, err := elasticsearch.NewClient(cfg)
			if err != nil {
				println("es client connect error")
				fmt.Println(err)
			}
			ESClient = &ES{
				Client: client,
			}
		},
	)
}
