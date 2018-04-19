package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aliyun-fc/ali_mns"
)

type appConf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func main() {
	conf := appConf{}

	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
		panic(e)
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			panic(e)
		}
	}

	client := ali_mns.NewAliMNSClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret)

	fcEndpoint := "<your account id>.cn-beijing.fc.aliyuncs.com/2016-08-15/services/ots-trigger/functions/mns-fc-test/invocations"
	topic := ali_mns.NewMNSTopic("test-fc-topic", client)

	sub := ali_mns.MessageSubsribeRequest{
		Endpoint:            topic.GenerateExtendEndpoint("fc", fcEndpoint),
		EndpointServiceName: "fc",
		EndpointRoleName:    "AliyunMNSPushFCRole",
		EndpointCipherSet:   "HmacSHA256",
		NotifyContentFormat: ali_mns.STREAM,
	}

	var err error
	err = topic.Subscribe("test-fc-subscibe", sub)
	if err != nil && !ali_mns.ERR_MNS_SUBSCRIPTION_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		fmt.Println(err)
		return
	}

	// now publish message
	msg := ali_mns.MessagePublishRequest{
		MessageBody: "hello topic <\"souriki/ali_mns\">",
		MessageAttributes: &ali_mns.MessageAttributes{
			//MailAttributes: &ali_mns.MailAttributes{
			//	Subject: "mns 触发 fc 测试",
			//	AccountName: "ls_huster@163.com",
			//},
			ExtendAttributes: &ali_mns.ExtendAttributes{
				Context: "fc only context",
			},
		},
	}
	_, err = topic.PublishMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
