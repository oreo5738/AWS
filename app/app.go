package app

import (
	"awsSDK/app/controller"
	"awsSDK/app/module"
	"context"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"strconv"
	"time"

	//"github.com/elastic/go-elasticsearch/v7/esapi"
	"os"
	//"strings"
)
type logset struct{
	Time  time.Time
	Event string
	Timestamp int64

}
func TransportData()  {
	//创建cloudwatch client
	fmt.Println(module.Configs.Esendpoint)
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(module.Configs.Region))
	cwcfg := cloudwatchlogs.NewFromConfig(cfg)
	/*获取日志组
	在conf中定义
	module.Configs.Log  要获取的日志组名
	module.Configs.Stream  日志流名
	*/
	timestamp ,err := controller.ESClient.SearchIndex()
	mygroup, err := cwcfg.GetLogEvents(context.TODO(),
		&cloudwatchlogs.GetLogEventsInput{
			StartTime: &timestamp,
			LogGroupName:  module.Configs.Log,
			LogStreamName: module.Configs.Stream,
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(timestamp)
	/*
		t 将日志中的时间戳取出来格式化成 2006-01-02 03:04:05
	*/
	for _, logs := range mygroup.Events {
		t, err := dateparse.ParseLocal(strconv.FormatInt(*logs.Timestamp, 10))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//声明一个实例，类型为logset，字段Time是格式化后的时间，字段Event是日志主体
		logsets := logset{
			Time : t,
			Event: *logs.Message,
			Timestamp: *logs.Timestamp,
		}
		//将logsets格式化成json,再将其转换成字符串
		jsonlog,_ := json.Marshal(logsets)
		jsonlogstr := string(jsonlog)

		/*
			与ES建立连接，发送请求，并关闭会话
			index : 在ES中创建的索引名
			Body ： 是要传入ES的json字符串
		*/
		controller.ESClient.RequestIndex(jsonlogstr)


	}

}
//测试函数
func Test()(Timestamp int64,err error )  {
	return  controller.ESClient.SearchIndex()
}