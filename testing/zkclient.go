package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	rootPath     = "/xxxx/xxxx/"
	categoryPath = "/xxxx/"
	outputPaht   = "/xxxx-0000000000"
)

type ZkConfig struct {
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Level    string    `json:"level"`
	Additivity string  `json:"additivity"`
	Tpls     []string  `json:"tpls"`
	Receiver string    `json:"receiver"`
	Qtype    string    `json:"type"`
	Attrs    AttrsData `json:"attrs"`
}
type AttrsData struct {
	Layout            string `json:"Layout"`
	ConversionPattern string `json:"ConversionPattern"`
	RecordInterval    int    `json:"RecordInterval"`
	Zookeeper         string `json:"Zookeeper"`
	Conf              string `json:"Conf"`
	Topic             string `json:"Topic"`
	CachedTime        string `json:"CachedTime"`
}


type ZKClient struct {
	zkServers []string
	conn      *zk.Conn
}

func main() {

	fi, err := os.Open("/topic.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	var zkserver []string
	zkserver = append(zkserver, "xxxx.xxxx.xxx.xx:port")
	zkClient, err := NewClient(zkserver, 5)
	if err != nil {
		fmt.Println("zkClient connect error:%v", err)
	}
	br := bufio.NewReader(fi)
	for {
		topic, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		t := strings.TrimSpace(string(topic))

		//time.Sleep(time.Second)
		err = zkClient.Create(t)
		if err != nil {
			fmt.Println("zkCreate topic error:%v%v", t, err)
		}

		/*err = zkClient.DeleteDir(t)
		if err != nil {
			fmt.Println("zkDelete topic error:", t, err)
		}*/

		/*err = zkClient.Get(t)
		if err != nil {
			fmt.Println("zkGet topic error:%v%v", t, err)
		}*/
	}
}



func NewClient(zkServers []string, timeout int) (*ZKClient, error) {
	client := new(ZKClient)
	client.zkServers = zkServers
	conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	client.conn = conn
	return client, nil
}


func (q *ZKClient) Get(topic string) error {
	path := rootPath + "shbt" + categoryPath + topic + outputPaht
	//path := "/qlog/search/zzzc/category/qlog_m-so-com-quc"
	v, s, err := q.conn.Get(path)
	if err != nil {
		q.conn.Close()
		return err
	}
	fmt.Println("v:%v",string(v),"s:%v",s.Aversion)

	return nil
}

func (q *ZKClient) Create(topic string) error {
	path := rootPath + "xxxx" + categoryPath + topic + outputPaht

	err := q.CreateDir(filepath.Dir(path), topic)
	if err != nil {
		return err
	}
	err = q.CreateDir(path, topic)
	if err != nil {
		return err
	}
	return nil
}
func (q *ZKClient) CreateDir(dir string, topic string) error {

	exists, _, err := q.conn.Exists(dir)
	if err != nil {
		q.conn.Close()
		return err
	}
	if !exists {
		data, err := json.Marshal(q.defaultZkConfig(topic))
		_, err = q.conn.Create(dir, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			q.conn.Close()
			return err
		}
		fmt.Println("dir:%v", dir, "create topic sucessed!")
	}

	return nil
}

func (q *ZKClient) DeleteDir(topic string) error {

	path := rootPath + "xxxx" + categoryPath + topic + outputPaht

	exists, s, err := q.conn.Exists(path)
	if err != nil {
		fmt.Println(err)
		q.conn.Close()
		return err
	}
	if exists {
		err := q.conn.Delete(path, s.Version)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("toplic:",topic,"path:",path,"delete sucess!")
	}

	parentPath := rootPath + "xxxx" + categoryPath + topic
	parentExists, s, err := q.conn.Exists(parentPath)
	if err != nil {
		fmt.Println(err)
		q.conn.Close()
		return err
	}
	if parentExists {
		err := q.conn.Delete(parentPath, s.Version)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("toplic:",topic,"parentPath:",path,"delete sucess!")
	}
	return nil
}

func (q *ZKClient) Close() {
	q.conn.Close()
}

func (q *ZKClient) defaultZkConfig(topic string) ZkConfig {
	return ZkConfig{
		Name: topic,
		Desc: "云平台日志",
		Level: "INFO",
		Additivity: "false",
		Receiver: "xxxx",
		Qtype:    "xxxx",
		Attrs: AttrsData{
			Layout:            "PatternLayout",
			ConversionPattern: "%D:%d{%q}\t[%-5p]\t[%5P:%14t]\t<%F:%L>\t%c\t%x\t%m%n",
			RecordInterval:    1,
			Zookeeper:         "xxxx",
			Conf:              "/producer.config",
			Topic:             topic,
			CachedTime:        "true",
		},
	}
}
