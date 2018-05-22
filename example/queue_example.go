package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gogap/logs"
	"github.com/souriki/ali_mns"
)

const (
	// MNSTestTimeout : run `MNS_TEST_TIMEOUT=true go run queue_example.go` to enable timeout test
	// should get lots of `timeout` in log.
	MNSTestTimeout = "MNS_TEST_TIMEOUT"
)

type appConf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	isTestTimeout := os.Getenv(MNSTestTimeout) != ""

	conf := appConf{}

	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
		panic(e)
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			panic(e)
		}
	}

	var client ali_mns.MNSClient
	if isTestTimeout {
		client = ali_mns.NewAliMNSClient(conf.Url,
			conf.AccessKeyId,
			conf.AccessKeySecret,
			ali_mns.Timeout(3))
	} else {
		client = ali_mns.NewAliMNSClient(conf.Url, conf.AccessKeyId, conf.AccessKeySecret)
	}

	msg := ali_mns.MessageSendRequest{
		MessageBody:  "hello <\"souriki/ali_mns\">",
		DelaySeconds: 0,
		Priority:     8}

	queueManager := ali_mns.NewMNSQueueManager(client)

	err := queueManager.CreateQueue("test", 0, 65536, 345600, 30, 0, 3)

	if err != nil && !ali_mns.ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		fmt.Println(err)
		return
	}

	queue := ali_mns.NewMNSQueue("test", client)

	for i := 1; i < 10000; i++ {
		var err error
		if isTestTimeout {
			_, err = queue.SendMessage(msg, ali_mns.RequestTimeout(10*time.Millisecond))
		} else {
			_, err = queue.SendMessage(msg)
		}

		go func() {
			fmt.Println(queue.QPSMonitor().QPS())
		}()

		if err != nil {
			fmt.Println(err)
		} else {
			// logs.Pretty("response:", ret)
		}

		endChan := make(chan int)
		respChan := make(chan ali_mns.MessageReceiveResponse)
		errChan := make(chan error)
		go func() {
			select {
			case resp := <-respChan:
				{
					// logs.Pretty("response:", resp)
					logs.Debug("change the visibility: ", resp.ReceiptHandle)
					if ret, e := queue.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
						fmt.Println(e)
					} else {
						// logs.Pretty("visibility changed", ret)
						logs.Debug("delete it now: ", ret.ReceiptHandle)
						if e := queue.DeleteMessage(ret.ReceiptHandle); e != nil {
							fmt.Println(e)
						}
						endChan <- 1
					}
				}
			case err := <-errChan:
				{
					fmt.Println(err)
					endChan <- 1
				}
			}
		}()

		queue.ReceiveMessage(respChan, errChan, 30)
		<-endChan
	}
}
