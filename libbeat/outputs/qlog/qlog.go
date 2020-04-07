package qlog

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"github.com/elastic/beats/libbeat/outputs/outil"
)

func init() {
	outputs.RegisterType("qlog", makeQlog)
}

func makeQlog(_ outputs.IndexManager, beat beat.Info, observer outputs.Observer, cfg *common.Config) (outputs.Group, error) {

	err := cfg.Unpack(&config)
	if err != nil {
		return outputs.Fail(err)
	}

	config, err := readConfig(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	topic, err := outil.BuildSelectorFromConfig(cfg, outil.Settings{
		Key:              "topic",
		MultiKey:         "topics",
		EnableSingleOnly: true,
		FailEmpty:        true,
	})
	if err != nil {
		return outputs.Fail(err)
	}

	hosts, err := outputs.ReadHostList(cfg)
	if err != nil {
		return outputs.Fail(err)
	}
	if !config.Container {
		zkClient, err := NewClient(config.ZkServer, 5)
		if err != nil {
			return outputs.Fail(err)
		}
		err = zkClient.Create(config)
		if err != nil {
			return outputs.Fail(err)
		}
	}

	InitConnectPool(hosts, time.Duration(config.TimeOut), config.MinConnNum)

	codec, _ := codec.CreateEncoder(beat, config.Codec)

	client, err := newQlogClient(observer, hosts, beat.IndexPrefix, topic, codec)
	if err != nil {
		return outputs.Fail(err)
	}

	retry := 0
	if config.MaxRetries < 0 {
		retry = -1
	}

	//将qlog 单条发送按照队列批量发送。设置Bulk最大值
	return outputs.Success(config.BulkMaxSize, retry, client)
}

func readConfig(cfg *common.Config) (*Config, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
