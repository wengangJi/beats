package qlog

import (
	"time"
	"path/filepath"
	"encoding/json"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	rootPath     = "/qlog/k8s/"
	categoryPath = "/category/"
	outputPaht   = "/output-0000000000"
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

func (q *ZKClient) Create(c *Config) error {
	path := rootPath + c.IDC + categoryPath + c.Topic + outputPaht

	err := q.CreateDir(filepath.Dir(path))
	if err != nil {
		return err
	}
	q.CreateDir(path)
	if err != nil {
		return err
	}
	return nil
}
func (q *ZKClient) CreateDir(dir string) error {

	exists, _, err := q.conn.Exists(dir)
	if err != nil {
		Close()
		return err
	}
	if !exists {
		data, err := json.Marshal(q.defaultZkConfig())
		_, err = q.conn.Create(dir, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			Close()
			return err
		}
	}
	return nil
}

func (q *ZKClient) Close() {
	q.conn.Close()
}



func (q *ZKClient) defaultZkConfig() ZkConfig {
	return ZkConfig{
		Name: config.Topic,
		Desc: "云平台日志",
		Level: "INFO",
		Additivity: "false",
		Receiver: "k8s_qlogd_"+config.IDC,
		Qtype:    "Qbus2",
		Attrs: AttrsData{
			Layout:            "PatternLayout",
			ConversionPattern: "%D:%d{%q}\t[%-5p]\t[%5P:%14t]\t<%F:%L>\t%c\t%x\t%m%n",
			RecordInterval:    1,
			Zookeeper:         config.Cluster,
			Conf:              "/home/s/apps/ucs_qlogd/etc/producer.config",
			Topic:             config.Topic,
			CachedTime:        "true",
		},
	}
}
